# Local Developer Analytics

It's normally a lot of guesswork to figure out the bottlenecks in the 
inner loop of the software dev lifecycle (SDLC). This project aims
to provide a bit of insight into the impact of day-to-day tooling
on a developer's workflows. 
 

# Open Considerations

## Deployment Model
So far, if we expect this to run on developer laptops, the interaction is either manually initiated or needs device management.

## What specific metrics do we want
* Time it takes to execute a command

## Product Requirements
* Aggregation endpoint that stores metadata and enables dashboards/insights to be built/rendered
* Dashboards that automatically highlight critical chokepoints/inefficiencies for developers -- things like command time execution specifically around tests running or waiting for rebuilds.
