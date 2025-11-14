# Implementation Plan

- [x] 1. Add service method to UserService
  - Create `UpcomingReminderResponse` struct in `internal/services/user_service.go`
  - Implement `GetUpcomingReminders(userID string)` method that:
    - Queries all families the user belongs to via `family_members` table
    - Retrieves memorial IDs from `memorial_families` for those families
    - Queries `memorial_reminders` for active reminders within 30 days
    - Enriches response with memorial and family information
    - Calculates `days_until` for each reminder
    - Returns sorted list by reminder_date ascending
  - _Requirements: 1.1, 1.2, 1.3, 1.4, 1.5, 2.1, 2.2, 2.3, 2.4, 2.5_

- [x] 2. Add controller method to UserController
  - Add `GetUpcomingReminders(ctx *gin.Context)` method in `internal/controllers/user_controller.go`
  - Extract user ID from JWT context
  - Call `userService.GetUpcomingReminders(userID)`
  - Return JSON response with code 0 on success
  - Handle errors with appropriate HTTP status codes and error messages
  - _Requirements: 1.1, 3.5_

- [x] 3. Register route in router
  - Add route `users.GET("/reminders/upcoming", userController.GetUpcomingReminders)` in `internal/router/router.go`
  - Ensure route is in the protected group (requires JWT authentication)
  - _Requirements: 3.1_

- [x] 4. Verify database indexes
  - Check that indexes exist on `family_members.user_id`
  - Check that indexes exist on `memorial_families.family_id`
  - Check that indexes exist on `memorial_reminders.memorial_id` and `memorial_reminders.reminder_date`
  - Document index status or create migration if needed
  - _Requirements: 3.2, 3.3_

- [ ]* 5. Add unit tests for service method
  - Test case: User with no families returns empty array
  - Test case: User with families but no reminders returns empty array
  - Test case: User with valid reminders returns correct data structure
  - Test case: Only reminders within 30 days are returned
  - Test case: Only active reminders are returned
  - Test case: Reminders are sorted by date ascending
  - _Requirements: 1.1, 1.2, 1.3, 1.4, 1.5, 2.1, 2.2, 2.3, 2.4, 2.5_

- [ ]* 6. Add integration test for API endpoint
  - Create test user, families, memorials, and reminders
  - Make authenticated GET request to `/api/v1/reminders/upcoming`
  - Verify response structure matches design
  - Verify response data is correct
  - Test unauthenticated request returns 401
  - _Requirements: 3.1, 3.4_

- [x] 7. Test with miniprogram
  - Start backend server
  - Open miniprogram and navigate to home page
  - Verify `/api/v1/reminders/upcoming` request succeeds (no 404)
  - Verify upcoming reminders are displayed correctly on home page
  - _Requirements: 1.1, 2.1, 2.2, 2.3, 2.4, 2.5_
