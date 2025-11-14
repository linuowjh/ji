# Requirements Document

## Introduction

This feature addresses a 404 error in the miniprogram where it requests `/api/v1/reminders/upcoming` to display upcoming memorial day reminders on the home page. Currently, the backend only provides family-specific reminder endpoints (`/api/v1/families/:family_id/reminders/upcoming`), but the miniprogram needs a user-level endpoint that aggregates reminders across all families the user belongs to.

## Glossary

- **System**: The Yun Nian Memorial backend API service
- **User**: An authenticated user of the miniprogram
- **Reminder**: A scheduled notification for a memorial day (e.g., birthday, death anniversary)
- **Family**: A family group that can have multiple members and associated memorials
- **Memorial Day**: A significant date associated with a deceased person (birthday, death anniversary, etc.)

## Requirements

### Requirement 1

**User Story:** As a miniprogram user, I want to see upcoming memorial day reminders from all my families on the home page, so that I can be aware of important dates without navigating to each family separately.

#### Acceptance Criteria

1. WHEN the User requests upcoming reminders, THE System SHALL return reminders from all families the User belongs to
2. THE System SHALL sort reminders by date in ascending order (nearest dates first)
3. THE System SHALL include family information with each reminder for context
4. THE System SHALL limit results to reminders within the next 30 days
5. THE System SHALL return an empty array when the User has no upcoming reminders

### Requirement 2

**User Story:** As a miniprogram user, I want the reminder data to include all necessary information, so that I can understand what the reminder is for without additional API calls.

#### Acceptance Criteria

1. THE System SHALL include the reminder date in each response item
2. THE System SHALL include the reminder type (birthday, death anniversary, etc.) in each response item
3. THE System SHALL include the associated memorial information (name, photo) in each response item
4. THE System SHALL include the family name in each response item
5. THE System SHALL include the days until the reminder date in each response item

### Requirement 3

**User Story:** As a system administrator, I want the endpoint to be performant and secure, so that it can handle multiple concurrent requests without degrading user experience.

#### Acceptance Criteria

1. THE System SHALL require JWT authentication for the endpoint
2. THE System SHALL complete the request within 500 milliseconds under normal load
3. THE System SHALL use database query optimization (joins, indexes) to minimize query time
4. THE System SHALL handle cases where the User belongs to no families gracefully
5. THE System SHALL log errors without exposing sensitive information to the client
