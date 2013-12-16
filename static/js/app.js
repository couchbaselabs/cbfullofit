'use strict';


// Declare app level module which depends on filters, and services
angular.module('myApp', [
  'ngRoute',
  'myApp.filters',
  'myApp.services',
  'myApp.directives',
  'myApp.controllers',
]).
config(['$routeProvider','$locationProvider', function($routeProvider, $locationProvider) {
  $routeProvider.when('/home', {templateUrl: '/static/partials/home.html', controller: 'HomeCtrl'});
  $routeProvider.when('/node', {templateUrl: '/static/partials/nodes.html', controller: 'NodesCtrl'});
  $routeProvider.when('/index', {templateUrl: '/static/partials/newindex.html', controller: 'NewIndexCtrl'});
  $routeProvider.when('/index/:indexName/schema', {templateUrl: '/static/partials/index-schema.html', controller: 'IndexCtrl'});
  $routeProvider.when('/index/:indexName/indexers', {templateUrl: '/static/partials/index-indexers.html', controller: 'IndexCtrl'});
  $routeProvider.when('/index/:indexName/search', {templateUrl: '/static/partials/index-search.html', controller: 'IndexCtrl'});
  $routeProvider.when('/index/:indexName/booleanSearch', {templateUrl: '/static/partials/index-boolean-search.html', controller: 'IndexCtrl'});
  $routeProvider.otherwise({redirectTo: '/home'});
  $locationProvider.html5Mode(true);
}]);
