# Daily Coding Problem: Problem #630 [Medium]

This problem was asked by Apple.

Implement a job scheduler which takes in a function f and an integer n,
and calls f after n milliseconds.

## Build and Run

I wrote two test programs to get the heap manipulation code correct.

If you have [GraphViz](https://graphviz.org/) installed, you can have the code build
a heap, and see what it looks like.

```sh
$ go build drawheap.go
$ ./drawheap 0 10 1 3 -2 6 9 8 11 4 > heap.dot
$ dot -Tpng -o heap.png heap.dot
```
You should check that the smallest numerical value is at the root,
no children are of greater value than their parent,
and that the bottom "row" of leaf nodes is filled in left-to-right

I wrote a program to sort integers via a heap.

```sh
$ go build sort.go
$ ./sort 0 10 1 3 -2 6 9 8 11 4

-2
0
1
3
4
6
8
9
10
11
```

Typical mutex-locking scheduler:

```sh
$ go build sched.go
$ ./sched 5000 1000 2000
Scheduling for wakeup at 2020-08-27T12:09:54.352041426-06:00

sleeping for 1s

Rescheduling for wakeup at 2020-08-27T12:09:52.352427131-06:00

1 Now:             2020-08-27T12:09:52.352765571-06:00
1 Wanted to run at 2020-08-27T12:09:52.352420691-06:00

Scheduling for wakeup at 2020-08-27T12:09:54.352041426-06:00

0 Now:             2020-08-27T12:09:54.352350527-06:00
0 Wanted to run at 2020-08-27T12:09:54.352040001-06:00

All scheduled jobs done
```

The command line means:

* Schedule a function to run in 5000 milliseconds (5 seconds)
* Sleep 1000 milliseconds
* Schedule a function to run in 2000 milliseconds

The first function gets scheduled to run in 5 seconds,
then the code sleeps 1 second.
That's 4 seconds until the first function should run.
After the sleep the code schedules a 2nd function to run in 2 seconds.
That's 2 seconds before the first function should run.

The functions scheduled to run have a serial number,
and a timestamp for when they should run.
When executed, the functions print their serial number,
when they ran,
and when they should have gotten run.

We can see that function serial number 0 (first function)
gets scheduled.
Then the second function (serial number 1) get scheduled to
run before s/n 0.
That actually happens,
and the desired execution time is
within 1 millsecond of the actual time a function executes.

By juggling scheduled times and sleep times,
you can try to get the code to reschedule execution times,
have 2 or more functions to run at the same wall clock time,
etc.

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

I chose to write a full-fledged job scheduler,
with a scheduling thread that runs in the background,
each function running in its own thread.
This just seemed more fun.

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
This is quite a lot to think of and do.
The choice of data structure for scheduling (a heap as a priority queue)
might be standard,
but the concurrency isn't.
Scheduling a job to run before the currently-scheduled-job should be run
isn't easy to get correct.

Job candidates could distinguish themselves not only by
actually writing (whiteboard!?) code,
but by noting difficult spots while talking through a design,
noting alternatives and why to not use them,
and also the usual "how to test",
and what test cases should occur.
Candidates who are versed in more than 1 operating system
could note different choices for each OS in scheduling primitives,
and concurrency primitives.
