Versioning
==========

Charts versioning
-----------------

Charts versioning is based on Semantic Versioning model. Version has a format
MAJOR.MINOR.PATCH.

For the version in ``Chart.yaml`` you should:

* increment MAJOR version if backwards incompatible changes are made. For
  example, some variables are changed.
* increment MINOR version if if you add new functionality, but those changes
  are backwards compatible. For example, if you introduce new variables, that
  are optional.
* increment PATCH version if no new functionality added, only backwards
  compatible fixes are made.

For the version in ``requirements.yaml`` you'll have to change MAJOR
version only.

Images versioning
-----------------

TODO
+ how image version change affects chart version