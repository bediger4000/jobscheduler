# Daily Coding Problem: Problem #630 [Medium]

This problem was asked by Apple.

Implement a job scheduler which takes in a function f and an integer n,
and calls f after n milliseconds.

## Analysis

This is a vague problem statement.
A job candidate could implement a one-shot timer and maybe meet the
carefully-parsed problem statement.

It seems likely that the problem requires a full-fledged job scheduler,
where new jobs can be added at any time.
This leads to interesting cases where the currently most urgent
job is schedule to run at time X,
but the newly to-be-scheduled job runs at time X - 1.
The "wakeup and run a function"
timer has to be reset to accomdate the to-be-scheduled job.

It's also possible that the interviewer would use this problem
for candidates of different nominal experience level,
expecting more from candidates with more experience.

There's a whole lot to this problem.

* Data structure to hold pending jobs.
I used a [binary heap](https://en.wikipedia.org/wiki/Binary_heap),
but others could work.
* How are pending jobs ordered?
Choice of data structure to hold jobs drives this.
* What OS primitive to use for scheduling?
* The problem statement (albeit vague) probably requires concurrency.
Decisions have to be made about mutex locks,
or other forms of concurrency.

I think that the "medium" rating is ridiculous.

Job candidates could distinguish themselves not only by
actually writing (whiteboard!?) code,
but by noting difficult spots while talking through a design,
noting alternatives and why to not use them,
and also the usual "how to test",
and what test cases should occur.

