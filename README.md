# Daily Coding Problem: Problem #630 [Medium]

This problem was asked by Apple.

Implement a job scheduler which takes in a function f and an integer n,
and calls f after n milliseconds.

## Analysis

This is a vague problem statement.
A job candidate could implement a one-shot timer and maybe meet the
carefully-parsed problem statement.

It seems likely that a more full-fledged job scheduler is desired,
where new jobs can be added at any time,
which can lead to interesting cases where the currently most urgent
job is schedule to run at time X,
but the newly to-be-scheduled job runs at time X - 1.
So the timer has to be reset to accomdate the to-be-scheduled job.

It's also possible that the interviewer would use this problem
for candidates of different nominal experience level,
expecting more from candidates with more experience.
