# MYN Backend Controller Verification Log

Tracks which Spring Boot controllers have been cross-referenced against the MYN-TUI-CLI spec,
and what was found. Each entry records what was verified, the commit where corrections were
applied, and any remaining known gaps.

---

## Legend

| Status | Meaning |
|--------|---------|
| ✅ Verified | Controller fully cross-checked; spec updated |
| ⚠️ Partial | Verified but some DTOs/fields still undocumented |
| ❌ Not verified | Not yet reviewed against spec |

---

## Verified Controllers

### ✅ UnifiedTaskController.java
**Path:** `/api/v2/unified-tasks`
**Verified in commit:** `6e446fa` (pass 4), `<pass5>` (pass 5)
**Endpoints verified:**
- GET / — list tasks ✓
- GET /{id} — get single task ✓
- POST / — create task ✓
- `@RequestMapping(value = "/{id}", method = {PUT, PATCH})` — update task (both PUT and PATCH accepted, MIN-169) ✓
- DELETE /{id} — delete task ✓
- POST /{id}/complete — complete task ✓
- POST /{id}/uncomplete — uncomplete task ✓
- POST /{id}/skip — skip habit ✓
- GET /{id}/streak — habit streak ✓
- POST /{id}/archive — archive task ✓
- GET /archived — list archived ✓
- PUT /{id}/unarchive — unarchive ✓
- POST /{id}/restore — restore deleted ✓
- PATCH /{id}/lock — lock task ✓
- PUT /{id}/convert — convert task type ✓
- PATCH /batch — batch update ✓
- DELETE /{id}/permanent — permanent delete ✓
- GET /missed — missed tasks ✓
- GET /completions — task completions history ✓
- GET /uuid/{taskUUID} — get by UUID ✓
- GET /schedule — scheduled tasks ✓
- POST /{id}/schedules/{scheduleId} — link schedule block ✓
- DELETE /{id}/schedules/{scheduleId} — unlink schedule block ✓
- GET /schedules/friendly — schedule-friendly task list ✓
- PUT /{id}/assignment-type — update assignment type ✓
- PUT /{id}/calendar/add — add to calendar ✓
- PUT /{id}/calendar/update — update calendar event ✓
- PUT /{id}/calendar/preferences — set calendar preferences ✓
- POST /{id}/calendar-complete — complete via calendar ✓
- POST /{id}/calendar-completion-choice — resolve completion prompt ✓
- GET /{id}/calendar/options — calendar options ✓
- GET /{id}/completion-status — completion status ✓
- GET /chores/today — chores today (alias in UnifiedTask controller) ✓
- POST /chores/instances/{instanceId}/complete — complete chore instance ✓
- GET /chores/statistics — chore stats ✓
- POST /households/{householdId}/chores — create household chore ✓
- GET /households/{householdId}/chores — list household chores ✓
- PUT /households/{householdId}/chores/{choreId} — update household chore ✓
- DELETE /households/{householdId}/chores/{choreId} — delete household chore ✓

**DTO fields verified:**
- `CreateTaskRequest`: title, startDate, dueDate, notes/description, priority, category, quantity, duration, isAutoScheduled, recurrenceRule, isLocked, schedules, projectId, isSlidingWindow, calendarId, scheduleInDev, allowSplitChunks, minChunkDuration
- `UpdateTaskRequest`: same as above plus completedDate, schedulingState, schedulingReason
- `UnifiedTaskDTO`: id, title, description/notes, taskType, priority, category, startDate, dueDate, duration, isCompleted, isArchived, recurrenceRule, createdDate, lastUpdated, streakCount, commentCount, schedules, calendarEvents, householdId, createdBy, ownerId, ownerName, skipAllowance, skipsUsed, bestStreak, difficulty, icon, timeWindow, parentTaskId, reminderEnabled, reminderTime, sharingMode, scope, completedToday

**Known gaps in spec:**
- Spec Appendix I Create Task only documents minimal fields; many optional fields exist in DTO
- Spec says `PATCH /api/v2/unified-tasks/{id}` — actually both PUT and PATCH work (MIN-169)
- `quantity`, `category`, `projectId`, `calendarId`, `isAutoScheduled` not in spec's Appendix I create request example

---

### ✅ ChoresController.java
**Path:** `/api/v2/chores`
**Verified in commit:** `6e446fa` (pass 4), `<pass5>` (pass 5)
**Endpoints verified:**
- GET /today — list today's chores (params: date, timezone, householdId) ✓
- POST /instances/{instanceId}/complete — complete chore instance ✓
- GET /schedule/{date} — schedule for a date ✓
- GET /schedule/range — schedule for a date range ✓
- GET /statistics — chore stats ✓
- PUT /{choreId}/schedule — assign/update schedule (added pass 5) ✓

**DTO fields verified:**
- `ChoreInstanceDTO`: instanceId, dateIdentifier, instanceDate, completed, completedDate, completedBy, completedByName, parentChoreId, title, description, icon, difficulty (1-3 create / 1-5 in DTO), recurrenceRule, householdId, assignedMemberIds, assignedMemberNames, assignedToCurrentUser, autoAssigned, position, active, completionType, estimatedMinutes, points, tags, notes, lastModified
- `ChoreStatsDTO`: memberStats (memberId, memberName, assignedCount, completedCount, completionRate), championId, championName, totalAssignments, totalCompleted, overallCompletionRate
- `ChoresPageResponse`: instances, heroStats (completionCount, streakDays, bestStreak, lastCompletedDate, completionRate, weeklyCompletions, memberName, avatarUrl), assignments, currentHeroMemberId, generatedAt, choreDate
- `CreateChoreRequest`: name (required, 1-100), icon (max 50), difficulty (1-3)
- `UpdateChoreRequest`: name, icon, position, difficulty (1-3)

**Known gaps:**
- Spec Appendix I has no ChoreInstanceDTO or ChoresPageResponse structure documented
- difficulty range: CreateChoreRequest limits to 1-3; ChoreInstanceDTO documents 1-5 (inconsistency in backend)

---

### ✅ CalendarController.java + CalendarV2Controller.java
**Path:** `/api/calendar/...` and `/api/v2/calendar/...`
**Verified in commit:** `<pass5>`
**Endpoints verified (user-facing, CLI-relevant):**
- GET /api/v2/calendar/events — list events ✓
- GET /api/v2/calendar/events/{id} — get single event (added pass 5) ✓
- POST /api/v2/calendar/standalone-events — create event ✓
- POST /api/v2/calendar/events — create (non-standalone) calendar event ✓
- DELETE /api/v2/calendar/events/{id} — delete event ✓
- POST /api/v2/calendar/meetings/{eventId}/decline — decline meeting ✓
- POST /api/v2/calendar/meetings/{eventId}/skip — skip meeting ✓
- GET /api/v2/calendar/meetings/skipped — list skipped meetings ✓
- DELETE /api/v2/calendar/meetings/skip/{meetingId} — unskip meeting ✓
- GET /api/v2/calendar/completion-preferences — calendar completion prefs (MIN-195) ✓
- PUT /api/v2/calendar/completion-preferences — update prefs (MIN-195) ✓

**Not documented in spec (internal/admin/account management):**
- POST /api/calendar/makeAccountPrimary
- POST /api/calendar/toggleCalendarIsUsing
- DELETE /api/calendar/deleteAccount
- POST /api/calendar/refreshTokens
- GET /api/calendar/mindyournow-events
- POST /api/calendar/emergency-cleanup
- POST /api/v2/calendar/completion-service/mark-completed
- POST /api/v2/calendar/completion-service/mark-incomplete
- GET /api/v2/calendar/sync-status/{id}
- GET /api/v2/tasks/{taskId}/calendar-events
- GET /api/v2/calendar/meetings/{eventId}/metadata
- DELETE /api/v2/calendar/meetings/{eventId}
- GET /api/v2/calendar/events/status
- GET /api/v2/calendar/events/status-by-event-id

**DTO fields verified:**
- `UserCalendarCompletionPreferencesDTO` (MIN-195): earlyCompletionDefault, earlyCompletionThresholdHours, lateCompletionAutoCleanupDays, autoHideCompletedAfterDays, showCompletionBadges, dimCompletedEvents, uncompleteBehavior

---

### ✅ PomodoroController.java
**Path:** `/api/v1/pomodoro`
**Verified in commit:** `6e446fa` (pass 4), `<pass5>` (pass 5)
**Endpoints verified:**
- POST /start ✓
- POST /smart-start ✓
- GET /suggestions?availableMinutes=&maxSuggestions= ✓
- POST /pause ✓
- POST /resume ✓
- POST /stop ✓
- POST /complete ✓
- GET /current ✓
- GET /stats (params: startDate, endDate) ✓
- GET /settings ✓
- PUT /settings ✓
- GET /sessions (paginated, filters: status, sessionType, taskId, startDate, endDate) ✓
- GET /sessions/{sessionId} ✓
- PUT /sessions/{sessionId} — update notes/interrupt ✓
- DELETE /sessions/{sessionId} ✓
- GET /household/{householdId}/activity (params: days) — not yet in spec ⚠️
- GET /household/{householdId}/stats (params: days) — not yet in spec ⚠️

**DTO files verified:** StartPomodoroRequest, SmartStartPomodoroRequest, UpdatePomodoroRequest, UpdatePomodoroSettingsRequest, TaskSuggestionsResponse, PomodoroStatsRequest, PomodoroStatsResponse, PomodoroSessionDTO, PomodoroSettingsDTO, PomodoroResponse

---

### ✅ TimerController.java
**Path:** `/api/v2/timers`
**Verified in commit:** `6e446fa` (pass 4), `<pass5>` (pass 5)
**Endpoints verified:**
- GET / — list timers (params: status, includeCompleted) ✓
- GET /{id} — get single timer (added pass 5) ✓
- POST /countdown — create countdown ✓
- POST /alarm — create alarm ✓
- POST /{id}/pause ✓
- POST /{id}/resume ✓
- POST /{id}/cancel ✓
- POST /{id}/snooze ✓
- POST /{id}/complete ✓
- DELETE /completed — dismiss all completed ✓
- DELETE /cancelled — clear all cancelled ✓
- GET /count ✓

**DTO fields verified:**
- `CreateCountdownTimerRequest`: name (required), durationSeconds (required, positive int), sourceTaskId, sourceHabitId, completionSound
- `CreateAlarmTimerRequest`: name (required), alarmTime (ZonedDateTime, required), sourceTaskId, sourceHabitId, completionSound. **NOTE: no recurrence field** — alarm `--recurrence` flag in CLI spec is not yet implemented in backend.
- `SnoozeTimerRequest`: snoozeMinutes (Integer, optional, 1-60)
- `TimerDTO`: id (UUID), name, type (COUNTDOWN|ALARM|AI_CREATED), durationSeconds, remainingSeconds, alarmTime, status, createdBy (USER|AI_KAIA|SYSTEM), sourceTaskId, sourceHabitId, escalationLevel, maxSnoozes, snoozeCount, createdAt, completedAt, pausedAt, expirationTime, completionSound

---

### ✅ TimerPreferencesController.java
**Path:** `/api/v2/customers/me/timer-preferences`
**Verified in commit:** `<pass5>`
**Endpoints verified:**
- GET / ✓
- PATCH / (**PATCH not PUT** — spec was wrong) ✓

**DTO fields verified:**
- `TimerPreferencesDTO`: completionSound (enum: `default|alarm|bell|chime|silent` — 5 values), defaultSnoozeMinutes (1-60), autoDismissEnabled, autoDismissDelaySeconds (0-300), showFloatingWidget, hapticEnabled
- **NOTE:** completionSound has only 5 options here vs 8 in TimerDTO/CreateTimerRequest (gong, ding, urgent, none not supported in preferences — backend inconsistency)

---

### ✅ HabitChainController.java
**Path:** `/api/habits/chains`
**Verified in commit:** `6e446fa` (pass 4)
**All endpoints match spec.** See entries in Appendix A.

---

### ✅ HabitReminderController.java
**Path:** `/api/habits/reminders`
**Verified in commit:** `6e446fa` (pass 4)
**Endpoints verified:**
- POST /{habitId}/calculate-smart-time ✓
- POST /{habitId}/test — not in spec (internal testing) ⚠️
- POST /check — admin-only, not in spec ⚠️

---

### ✅ HabitSchedulingController.java
**Path:** `/api/v2/scheduling`
**Verified in commit:** `6e446fa` (pass 4), `<pass5>`
**Endpoints verified:**
- POST /habits/schedule ✓
- GET /habits/status — added pass 5 ✓

---

### ✅ CustomerController.java
**Path:** `/api/v1/customers`
**Verified in commit:** `6e446fa` (pass 4)
**All documented endpoints confirmed correct.** See Appendix A entries in spec.

---

### ✅ NotificationController.java
**Path:** `/api/v1/notifications`
**Verified in commit:** `<pass5>`
**CRITICAL: All notification endpoints are v1, not v2:**
- GET /api/v1/notifications — list (spec was saying v2) ✓ fixed
- GET /api/v1/notifications/unread-count — unread count (spec was saying /api/v2/notifications/unread) ✓ fixed
- PUT /api/v1/notifications/{id}/read — mark single read (spec was saying POST v2/mark-read) ✓ fixed
- PUT /api/v1/notifications/read-all — mark all read (spec was saying POST v2/mark-read) ✓ fixed
- DELETE /api/v1/notifications/{id} — delete (spec was saying v2) ✓ fixed

**Note:** `MarkNotificationsReadRequest` DTO (with notificationIds array) exists in codebase but is **unused** — it's likely residual from a planned v2 refactor.

---

### ✅ McpSessionsController.java
**Path:** `/api/v1/customers/mcp-sessions`
**Verified in commit:** `<pass5>`
**Endpoints verified:**
- GET / — list sessions ✓
- DELETE /{id} — revoke session (added pass 5) ✓
- DELETE / — revoke all sessions (added pass 5) ✓
- GET /count — session count ✓

---

### ✅ PlanningController.java
**Path:** `/planning` (no `/api/` prefix — internal/admin only)
**Verified in commit:** `<pass5>`
**Endpoints:**
- GET /planning/plan — internal version/build info
- GET /planning/scheduleAll — bulk auto-schedule (internal admin)
- GET /planning/unScheduleAll — bulk un-schedule (internal admin)
- GET /planning/kickTheCan?rebalance=bool — rebalance scheduled tasks (internal admin)

**These are internal admin endpoints, not public user API.** Not documented in spec intentionally.

---

## Not Yet Verified

| Controller | Priority | Notes |
|------------|----------|-------|
| ConversationController.java | Medium | AI conversations |
| AIStreamingController.java | Medium | AI chat stream |
| GamificationController.java | Low | Verified manually in pass 3 |
| HouseholdController.java + HouseholdInviteController.java | Low | Verified in passes 2-4 |
| ProjectController.java | Low | Verified in pass 3 |
| CompassController.java | Low | Verified in pass 3 |
| SearchController.java | Low | Verified in pass 3 |
| McpOAuthController.java | Low | Verified in pass 2 |
| AccountDeletionController.java | Low | Verified in passes 3-4 |
| ExportController (in CustomerController) | Low | Verified in pass 4 |
| ScheduleController.java | Low | Verified in pass 3 (block CRUD only) |
| GroceryListController.java | Low | Verified in pass 3 |
