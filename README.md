# Local Developer Analytics

It's normally a lot of guesswork to figure out the bottlenecks in the 
inner loop of the software dev lifecycle (SDLC). This project aims
to provide a bit of insight into the impact of day-to-day tooling
on a developer's workflows. 

This implements metrics around step 2 in DORA metrics, lead time to changes (LT).

# Open Considerations
## Deployment Model
So far, if we expect this to run on developer laptops, the interaction is either manually initiated or needs device management.

## What specific metrics do we want
* Time it takes to execute a command
* Time spent waiting on network calls (i.e. downloading packages for dev)

## Product Requirements
* Aggregation endpoint that stores metadata and enables dashboards/insights to be built/rendered
* Dashboards that automatically highlight critical chokepoints/inefficiencies for developers -- things like command time execution specifically around tests running or waiting for rebuilds.

## Project Phases

### POC
* determine viable approaches for tracking command execution on macs (dtrace, openBSM, shell wrapping)
* determine viable approaches for tracking command execution on linux (bpftrace, Linux Audit Framework, shell wrapping)
* delieverable: end to end POCs for chosen scenarios that demonstrate the core ability of tracing command execution times and network download times

Goal: understand approach and UX tradeoffs for end users
Goal: understand distinction between commands taking time for network downloads versus commands that take time compiling

# DORA
Wrong place to be storing this right now, but since it's semi related and important context, adding it here for now.

## Things To Measure
* Code Review Speed
** Pull Request Size: Large pull requests can be difficult to review and merge, leading to delays.
** Pull Request Review Time: Similarly, pull requests that remain open for a long time before being merged can indicate bottlenecks in the code review process.
* Overall Velocity
** Commit Velocity: Understanding of time-to-commits on main branch and overall weekly/monthly velocity over time.
* Engineering Excellence
** Test Coverage: Untested codes leads to bugs, which leads to more time spent fixing bugs. Understand total test coverage of code.