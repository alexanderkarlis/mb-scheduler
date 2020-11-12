Select * from 
    (Select *, count(*)
    OVER
        (PARTITION BY
            runtime, fullname, username, password, classtime, weekday, date
        ) AS count
        from schedule_rt) tableWithCount
    Where tableWithCount.count > 1;
