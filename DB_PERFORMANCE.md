# Database Performance Guide

This document tracks potential database performance optimizations for future implementation as the application scales.

## Current Indexes

The following indexes are automatically created via GORM tags:

| Table | Column(s) | Index Type |
|-------|-----------|------------|
| `workout_prescriptions` | `workout_id` | B-tree |
| `workout_prescriptions` | `exercise_id` | B-tree |
| `workout_prescriptions` | `group_id` | B-tree |
| `exercise_muscle_groups` | `exercise_id`, `muscle_group_id` | Composite Unique |
| `refresh_tokens` | `user_id` | B-tree |
| `refresh_tokens` | `expires_at` | B-tree |

## Workout Filter Performance

### Title Search (`ILIKE '%term%'`)

**Current behavior:** Full table scan - cannot use standard B-tree index for leading wildcard searches.

**Scaling threshold:** Performance degrades noticeably at ~10k+ workouts per user.

**Future optimization:**
```sql
-- Enable trigram extension (one-time)
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- Add GIN trigram index for fast text search
CREATE INDEX idx_workouts_title_trgm ON workouts USING GIN (title gin_trgm_ops);
```

**Alternative:** For even better search, consider PostgreSQL full-text search:
```sql
-- Add tsvector column and index
ALTER TABLE workouts ADD COLUMN title_search tsvector
  GENERATED ALWAYS AS (to_tsvector('english', title)) STORED;
CREATE INDEX idx_workouts_title_fts ON workouts USING GIN (title_search);
```

### Muscle Group Filter (Subquery with JOIN)

**Current behavior:** Uses indexed foreign keys but involves subquery overhead.

**Scaling threshold:** May slow down at ~100k+ workout_prescriptions rows.

**Future optimization - Option 1:** Add composite index
```sql
CREATE INDEX idx_wp_exercise_workout ON workout_prescriptions (exercise_id, workout_id)
  WHERE deleted_at IS NULL;
```

**Future optimization - Option 2:** Denormalize by caching muscle group IDs on workouts
```sql
-- Add array column to workouts table
ALTER TABLE workouts ADD COLUMN muscle_group_ids uuid[] DEFAULT '{}';
CREATE INDEX idx_workouts_muscle_groups ON workouts USING GIN (muscle_group_ids);

-- Update via trigger or application logic when prescriptions change
```

### Exercise Filter (Simple Subquery)

**Current behavior:** Direct index lookup on `exercise_id` - performant.

**Status:** No optimization needed - existing index is sufficient.

## General Recommendations

### Query Monitoring

Add slow query logging to PostgreSQL:
```sql
-- In postgresql.conf
log_min_duration_statement = 100  -- Log queries taking > 100ms
```

### Connection Pooling

For high traffic, consider PgBouncer or built-in GORM connection pool tuning:
```go
sqlDB, _ := db.DB()
sqlDB.SetMaxIdleConns(10)
sqlDB.SetMaxOpenConns(100)
sqlDB.SetConnMaxLifetime(time.Hour)
```

### Pagination Best Practices

Current implementation uses OFFSET/LIMIT which can be slow for deep pagination.

**Future optimization:** Keyset (cursor) pagination for large result sets:
```sql
-- Instead of: OFFSET 10000 LIMIT 20
-- Use: WHERE created_at < ? ORDER BY created_at DESC LIMIT 20
```

## Implementation Priority

| Optimization | Priority | Trigger |
|-------------|----------|---------|
| Trigram index for title search | Medium | When search queries exceed 100ms |
| Composite index for muscle group filter | Low | When filter queries exceed 100ms |
| Cursor-based pagination | Low | When users paginate beyond page 100 |
| Full-text search | Low | When advanced search features needed |

## Monitoring Queries

Check for slow queries:
```sql
SELECT query, calls, mean_time, total_time
FROM pg_stat_statements
ORDER BY mean_time DESC
LIMIT 10;
```

Check index usage:
```sql
SELECT relname, indexrelname, idx_scan, idx_tup_read
FROM pg_stat_user_indexes
WHERE schemaname = 'public'
ORDER BY idx_scan DESC;
```
