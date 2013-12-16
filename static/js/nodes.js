function NodesCtrl($scope, $http) {
	$http.get('/api/node/').
	success(function(data) {
		$scope.nodes = data;
	}).
	error(function(data, code) {

	});

}