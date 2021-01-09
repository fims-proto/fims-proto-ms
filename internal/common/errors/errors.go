package errors

/*
error handling
should have 2 types of errors:
  1. business error:
    - thrown only from domain/app layer, as domain/app layer only cares about business logic,
      but not technical implementation.

  2. technical error:
    - thrown only from port/adapter layer, as port/adapter layer only cares about how to implement,
      but not business logic.

in some cases, panic directly, e.g. dependency injection failure (initialization fail)



TBD: wrap errors here, a separated package, or in each packages. We'll see...
*/
