// Package jobs provides a set of task grouping tools for running sequences of tasks, halting on error,
// and displaying their results.
//
// The four main entrypoints are Run, RunContext, Parallel, and ParallelContext. They accept implementations
// of the Job interface, which you can either implement in your own task types, or simply cast functions to the
// JobFunc type. They return Events, a synchronous bus of events that, for each job, always occur in the following
// order:
//
//  *EventQueued
//  *EventStarted
//  *EventProgressed
//  *EventFinished
//
// Additionally, if a job spawns any sub-jobs, their events will arrive between *EventStarted and *EventFinished.
package jobs
