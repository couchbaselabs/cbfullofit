'use strict';

/* Filters */

angular.module('myApp.filters', ['ngSanitize']).
  filter('interpolate', ['version', function(version) {
    return function(text) {
      return String(text).replace(/\%VERSION\%/mg, version);
    };
  }]);
