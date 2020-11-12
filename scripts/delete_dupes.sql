WITH unique_values AS (
    SELECT DISTINCT ON (
            runtime,
            fullname,
            username,
            password,
            classtime,
            weekday,
            date
        ) *
    FROM schedule_rt
)
DELETE FROM schedule_rt
WHERE schedule_rt.id NOT IN (
        SELECT id
        FROM unique_values as y
    );