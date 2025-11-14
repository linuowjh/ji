# Design Document

## Overview

This design document outlines the implementation of a user-level upcoming reminders API endpoint that aggregates memorial day reminders from all families a user belongs to. The endpoint will be accessible at `/api/v1/reminders/upcoming` and will provide a unified view of upcoming memorial dates across all of the user's family groups.

## Architecture

### High-Level Flow

1. User makes authenticated GET request to `/api/v1/reminders/upcoming`
2. System extracts user ID from JWT token
3. System queries all families the user belongs to
4. System retrieves reminders from all associated memorials
5. System filters reminders to next 30 days
6. System enriches reminder data with family and memorial information
7. System calculates days until each reminder
8. System sorts by date (ascending) and returns response

### Components Involved

- **Router** (`internal/router/router.go`): Add new route definition
- **Controller** (`internal/controllers/user_controller.go`): Add new handler method
- **Service** (`internal/services/user_service.go`): Add business logic for aggregating reminders
- **Models** (`internal/models/family.go`): Use existing `MemorialReminder`, `FamilyMember`, `MemorialFamily` models

## Components and Interfaces

### 1. Router Configuration

Add new route in the protected users group:

```go
users.GET("/reminders/upcoming", userController.GetUpcomingReminders)
```

### 2. Controller Method

Add method to `UserController`:

```go
// GetUpcomingReminders 获取用户所有家族的即将到来的纪念日提醒
func (c *UserController) GetUpcomingReminders(ctx *gin.Context) {
    userID, exists := ctx.Get("user_id")
    if !exists {
        ctx.JSON(http.StatusUnauthorized, APIResponse{
            Code:    1002,
            Message: "用户未登录",
        })
        return
    }

    reminders, err := c.userService.GetUpcomingReminders(userID.(string))
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, APIResponse{
            Code:    1005,
            Message: err.Error(),
        })
        return
    }

    ctx.JSON(http.StatusOK, APIResponse{
        Code:    0,
        Message: "获取成功",
        Data:    reminders,
    })
}
```

### 3. Service Method

Add method to `UserService`:

```go
// UpcomingReminderResponse 即将到来的提醒响应结构
type UpcomingReminderResponse struct {
    ID           string    `json:"id"`
    ReminderType string    `json:"reminder_type"`
    ReminderDate time.Time `json:"reminder_date"`
    Title        string    `json:"title"`
    Content      string    `json:"content"`
    DaysUntil    int       `json:"days_until"`
    Memorial     struct {
        ID        string `json:"id"`
        Name      string `json:"name"`
        AvatarURL string `json:"avatar_url"`
    } `json:"memorial"`
    Family struct {
        ID   string `json:"id"`
        Name string `json:"name"`
    } `json:"family"`
}

// GetUpcomingReminders 获取用户所有家族的即将到来的纪念日提醒
func (s *UserService) GetUpcomingReminders(userID string) ([]*UpcomingReminderResponse, error) {
    // 1. Get all families the user belongs to
    var familyMembers []models.FamilyMember
    err := s.db.Where("user_id = ?", userID).Find(&familyMembers).Error
    if err != nil {
        return nil, err
    }

    if len(familyMembers) == 0 {
        return []*UpcomingReminderResponse{}, nil
    }

    // 2. Extract family IDs
    familyIDs := make([]string, len(familyMembers))
    for i, member := range familyMembers {
        familyIDs[i] = member.FamilyID
    }

    // 3. Get all memorial IDs associated with these families
    var memorialFamilies []models.MemorialFamily
    err = s.db.Where("family_id IN ?", familyIDs).Find(&memorialFamilies).Error
    if err != nil {
        return nil, err
    }

    if len(memorialFamilies) == 0 {
        return []*UpcomingReminderResponse{}, nil
    }

    // 4. Create a map of memorial ID to family IDs for later lookup
    memorialToFamilies := make(map[string][]string)
    for _, mf := range memorialFamilies {
        memorialToFamilies[mf.MemorialID] = append(memorialToFamilies[mf.MemorialID], mf.FamilyID)
    }

    memorialIDs := make([]string, 0, len(memorialToFamilies))
    for memorialID := range memorialToFamilies {
        memorialIDs = append(memorialIDs, memorialID)
    }

    // 5. Query reminders within next 30 days
    now := time.Now()
    thirtyDaysLater := now.AddDate(0, 0, 30)

    var reminders []models.MemorialReminder
    err = s.db.Preload("Memorial").
        Where("memorial_id IN ? AND is_active = ? AND reminder_date BETWEEN ? AND ?",
            memorialIDs, true, now.Format("2006-01-02"), thirtyDaysLater.Format("2006-01-02")).
        Order("reminder_date ASC").
        Find(&reminders).Error
    if err != nil {
        return nil, err
    }

    // 6. Load family information
    var families []models.Family
    err = s.db.Where("id IN ?", familyIDs).Find(&families).Error
    if err != nil {
        return nil, err
    }

    familyMap := make(map[string]*models.Family)
    for i := range families {
        familyMap[families[i].ID] = &families[i]
    }

    // 7. Build response with enriched data
    responses := make([]*UpcomingReminderResponse, 0, len(reminders))
    for _, reminder := range reminders {
        // Get the first family associated with this memorial
        familyIDsForMemorial := memorialToFamilies[reminder.MemorialID]
        if len(familyIDsForMemorial) == 0 {
            continue
        }

        family := familyMap[familyIDsForMemorial[0]]
        if family == nil {
            continue
        }

        // Calculate days until reminder
        daysUntil := int(reminder.ReminderDate.Sub(now).Hours() / 24)

        response := &UpcomingReminderResponse{
            ID:           reminder.ID,
            ReminderType: reminder.ReminderType,
            ReminderDate: reminder.ReminderDate,
            Title:        reminder.Title,
            Content:      reminder.Content,
            DaysUntil:    daysUntil,
        }

        response.Memorial.ID = reminder.Memorial.ID
        response.Memorial.Name = reminder.Memorial.Name
        response.Memorial.AvatarURL = reminder.Memorial.AvatarURL

        response.Family.ID = family.ID
        response.Family.Name = family.Name

        responses = append(responses, response)
    }

    return responses, nil
}
```

## Data Models

### Existing Models Used

1. **FamilyMember**: Links users to families
   - Fields: `ID`, `FamilyID`, `UserID`, `Role`, `JoinedAt`

2. **MemorialFamily**: Links memorials to families
   - Fields: `ID`, `MemorialID`, `FamilyID`, `CreatedAt`

3. **MemorialReminder**: Stores reminder information
   - Fields: `ID`, `MemorialID`, `ReminderType`, `ReminderDate`, `Title`, `Content`, `IsActive`

4. **Memorial**: Memorial information
   - Fields: `ID`, `Name`, `AvatarURL`, etc.

5. **Family**: Family information
   - Fields: `ID`, `Name`, `CreatorID`, etc.

### Response Structure

```json
{
  "code": 0,
  "message": "获取成功",
  "data": [
    {
      "id": "reminder-uuid",
      "reminder_type": "birthday",
      "reminder_date": "2025-11-20T00:00:00Z",
      "title": "张三的生日",
      "content": "今天是张三的生日纪念日",
      "days_until": 6,
      "memorial": {
        "id": "memorial-uuid",
        "name": "张三",
        "avatar_url": "https://example.com/avatar.jpg"
      },
      "family": {
        "id": "family-uuid",
        "name": "张家大院"
      }
    }
  ]
}
```

## Error Handling

### Error Scenarios

1. **User not authenticated**: Return 401 with code 1002
2. **Database query error**: Return 500 with code 1005
3. **User has no families**: Return empty array (not an error)
4. **No upcoming reminders**: Return empty array (not an error)

### Error Response Format

```json
{
  "code": 1005,
  "message": "数据库查询错误: [error details]"
}
```

## Testing Strategy

### Unit Tests

1. **Test GetUpcomingReminders with no families**
   - User belongs to no families
   - Expected: Empty array returned

2. **Test GetUpcomingReminders with families but no reminders**
   - User belongs to families
   - Families have memorials
   - No active reminders within 30 days
   - Expected: Empty array returned

3. **Test GetUpcomingReminders with valid reminders**
   - User belongs to multiple families
   - Families have memorials with reminders
   - Expected: Sorted list of reminders with correct data

4. **Test GetUpcomingReminders filters by date range**
   - Reminders exist beyond 30 days
   - Expected: Only reminders within 30 days returned

5. **Test GetUpcomingReminders filters inactive reminders**
   - Some reminders have `is_active = false`
   - Expected: Only active reminders returned

6. **Test GetUpcomingReminders sorts by date**
   - Multiple reminders with different dates
   - Expected: Reminders sorted by date ascending

### Integration Tests

1. **Test full API endpoint**
   - Create test user, families, memorials, and reminders
   - Make authenticated request to `/api/v1/reminders/upcoming`
   - Verify response structure and data

2. **Test authentication requirement**
   - Make request without JWT token
   - Expected: 401 Unauthorized

3. **Test performance with multiple families**
   - User belongs to 10+ families
   - Each family has multiple memorials
   - Verify response time < 500ms

## Performance Considerations

### Database Query Optimization

1. **Use IN clauses**: Batch queries for families and memorials
2. **Preload associations**: Use GORM's Preload for Memorial data
3. **Index usage**: Ensure indexes exist on:
   - `family_members.user_id`
   - `memorial_families.family_id`
   - `memorial_reminders.memorial_id`
   - `memorial_reminders.reminder_date`
   - `memorial_reminders.is_active`

### Query Execution Plan

```
1. Query family_members WHERE user_id = ? (indexed)
2. Query memorial_families WHERE family_id IN (...) (indexed)
3. Query memorial_reminders WHERE memorial_id IN (...) AND is_active = true AND reminder_date BETWEEN ? AND ? (indexed)
4. Query families WHERE id IN (...) (indexed, primary key)
5. Preload memorials (already loaded with reminders)
```

Expected total query time: < 100ms for typical user with 5 families

### Caching Strategy (Future Enhancement)

- Cache user's family memberships for 5 minutes
- Cache reminder data for 1 minute
- Invalidate cache when reminders are created/updated/deleted

## Security Considerations

1. **Authentication**: JWT token required (enforced by middleware)
2. **Authorization**: User can only see reminders from families they belong to
3. **Data exposure**: Only expose necessary memorial and family information
4. **SQL injection**: Use parameterized queries (GORM handles this)
5. **Rate limiting**: Existing rate limit middleware applies

## Migration Requirements

No database migrations required. All necessary tables and columns already exist:
- `family_members`
- `memorial_families`
- `memorial_reminders`
- `memorials`
- `families`

## Deployment Considerations

1. **Backward compatibility**: New endpoint, no breaking changes
2. **Rollback plan**: Simply remove the route if issues arise
3. **Monitoring**: Add logging for query performance
4. **Feature flag**: Not required (low-risk addition)

## Future Enhancements

1. **Pagination**: Add page/page_size parameters if users have many reminders
2. **Filtering**: Allow filtering by reminder_type
3. **Custom date range**: Allow specifying days ahead (default 30)
4. **Push notifications**: Integrate with notification service
5. **Reminder preferences**: Allow users to configure which reminders they want to see
