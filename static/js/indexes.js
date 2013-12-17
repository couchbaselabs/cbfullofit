function HomeCtrl($scope, $http) {
	var updateIndexList = function() {
		$http.get('/api/index/').
		success(function(data) {
			$scope.indexes = data;
		}).
		error(function(data, code) {

		});
	};

	$scope.indexes = [];
	updateIndexList();

	$scope.deleteIndex = function(name) {
		var r = window.confirm("Are you sure you want to delete the index named '" + name + "'?");
		if (r === true) {
			$http.delete('/api/index/' + name).
			success(function(data) {
				updateIndexList();
			}).
			error(function(data, code) {

			});
		}
	};
}

function NewIndexCtrl($scope, $http, $log) {

	var resetSchema = function() {
		$scope.fieldName = "";
		$scope.fieldPath = "";
		$scope.fieldAnalyzer = "";
	};

	var resetForm = function() {
		$scope.indexName = "";
		$scope.fields = {};
		resetSchema();
	};


	resetForm();
	$http.get('/api/bucket/').
	success(function(data) {
		$scope.buckets = data;
		$scope.indexBucket = $scope.buckets[0];
	}).
	error(function(data, code) {

	});

	$scope.removeField = function(name) {
		delete $scope.fields[name];
	};

	$scope.addField = function() {

		existing = $scope.fields[$scope.fieldName];
		if(existing !== undefined) {
			$scope.errorMessage = "Field name " + $scope.fieldName + " already in use, choose a different name";
			return;
		}
		
		if($scope.fieldName === "") {
			$scope.errorMessage = "Schema field name cannot be empty";
			return;
		}

		if($scope.fieldPath === "") {
			$scope.errorMessage = "Schema field path cannot be empty";
			return;
		}

		if($scope.fieldAnalyzer === "") {
			$scope.errorMessage = "Select a field analyzer";
			return;
		}

		field = {
			"path": $scope.fieldPath,
			"analyzer": $scope.fieldAnalyzer
		};

		$scope.fields[$scope.fieldName] = field;

		$log.debug($scope.fields);

		// reset form
		delete $scope.errorMessage;
		resetSchema();
	};

	$scope.createIndex = function() {

		if ($scope.indexName === "") {
			$scope.errorMessage = "Index name cannot be empty";
			return;
		}

		if (Object.keys($scope.fields).length === 0) {
			$scope.errorMessage = "Please add at least one field to index";
			return;
		}

		// try to create the index
		requestBody = {
			"type": "index",
			"name": $scope.indexName,
			"bucket": $scope.indexBucket,
			"schema": $scope.fields
		};

		$http.put('/api/index/' + $scope.indexName, requestBody).
		success(function(data) {
			delete $scope.errorMessage;
			resetForm();
			// redirect to new index page
		}).
		error(function(data, code) {
			$scope.errorMessage = "Unable to create index: " + data;
		});
	};

}

function IndexCtrl($scope, $http, $routeParams, $log) {

	$scope.minShouldOptions = [];
	for (var i = 0; i <= 50; i++) {
		$scope.minShouldOptions[i] = i;
	}

	var resetSchema = function() {
		$scope.clauseTerm = "";
		$scope.clauseOccur = "MUST";
		$scope.clauseBoost = 1.0;
		for(var f in $scope.theindex.schema) {
			$scope.clauseField = f;
			break;
		}
	};

	var resetForm = function() {
		$scope.clauses = [];
		$scope.size = "10";
		$scope.minShould = "0";
		resetSchema();
	};

	$http.get('/api/index/' + $routeParams.indexName).
	success(function(data) {
		$scope.theindex = data;
		for(var f in $scope.theindex.schema) {
			$scope.field = f;
			break;
		}
		$log.debug($scope.theindex);
		$log.debug($scope.theindex.schema);
		resetForm();
	}).
	error(function(data, code) {

	});

	$http.get('/api/node/').
	success(function(data) {
		$scope.nodes = data;
	}).
	error(function(data, code) {

	});

	$scope.nodesAssigned = {};
	$http.get('/api/assignment/index/' + $routeParams.indexName).
	success(function(data) {
		for(var i in data) {
			$scope.nodesAssigned[data[i]] = true;
		}
	}).
	error(function(data, code) {

	});

	$scope.toggleNode = function(node) {
		
		if ($scope.nodesAssigned[node]) {
			// do delete
			$http.delete('/api/assignment/index/' + $routeParams.indexName + '/' + node).
			success(function(data) {
				delete $scope.nodesAssigned[node];
			}).
			error(function(data, code) {

			});
		} else {
			// do assign
			$http.put('/api/assignment/index/' + $routeParams.indexName + '/' + node, '').
			success(function(data) {
				$scope.nodesAssigned[node] = true;
			}).
			error(function(data, code) {

			});
		}
	};

	$scope.search = function() {
		$http.get('/api/index/' + $scope.theindex.name + '/_searchTerm?q=' + $scope.term + '&f=' + $scope.field).
		success(function(data) {
			$scope.results = data;
			for(var i in $scope.results.hits) {
				hit = $scope.results.hits[i];
				hit.roundedScore = $scope.roundScore(hit.score);
				hit.explanationString = $scope.expl(hit.explanation);
			}
		}).
		error(function(data, code) {

		});
	};

	$scope.expl = function(explanation) {
		rv = "" + $scope.roundScore(explanation.value) + " - " + explanation.message;
		rv = rv + "<ul>";
		for(var i in explanation.children) {
			child = explanation.children[i];
			rv = rv + "<li>" + $scope.expl(child) + "</li>";
		}
		rv = rv + "</ul>";
		return rv;
	};

	$scope.roundScore = function(score) {
		return Math.round(score*1000)/1000;
	};

	$scope.removeClause = function(index) {
		console.log("remove");
		console.log(index);
		$scope.clauses.splice(index, 1);
	};

	$scope.addClause = function() {
		
		if($scope.clauseTerm === "") {
			$scope.errorMessage = "Clause term cannot be empty";
			return;
		}

		if($scope.clauseOccur === "") {
			$scope.errorMessage = "Select clause occur";
			return;
		}

		if($scope.clauseField === "") {
			$scope.errorMessage = "Select a field";
			return;
		}

		if($scope.clauseBoost === "") {
			$scope.errorMessage = "Clause boost cannot be empty";
			return;
		}

		clause = {
			"term": $scope.clauseTerm,
			"occur": $scope.clauseOccur,
			"field": $scope.clauseField,
			"boost": $scope.clauseBoost
		};

		$scope.clauses.push(clause);

		$log.debug($scope.clauses);

		// reset form
		delete $scope.errorMessage;
		resetSchema();
	};

	$scope.searchBoolean = function() {
		var requestBody = {
			"query": {
				"must": {
					"terms":[],
					"boost": 1.0,
					"explain": true
				},
				"should":{
					"terms":[],
					"boost": 1.0,
					"explain": true,
					"min": parseInt($scope.minShould, 10)
				},
				"must_not": {
					"terms": [],
					"boost": 1.0,
					"explain": true
				},
				"boost": 1.0,
				"explain": true
			},
			explain: true,
			size: parseInt($scope.size, 10)
		};
		for(var i in $scope.clauses) {
			var clause = $scope.clauses[i];
			var termQuery = {
				"term": clause.term,
				"field": clause.field,
				"boost": clause.boost,
				"explain": true
			};
			switch(clause.occur) {
				case "MUST":
				requestBody.query.must.terms.push(termQuery);
				break;
				case "SHOULD":
				requestBody.query.should.terms.push(termQuery);
				break;
				case "MUST NOT":
				requestBody.query.must_not.terms.push(termQuery);
				break;
			}
		}
		if (requestBody.query.must.terms.length === 0) {
			delete requestBody.query.must;
		}
		if (requestBody.query.should.terms.length === 0) {
			delete requestBody.query.should;
		}
		if (requestBody.query.must_not.terms.length === 0) {
			delete requestBody.query.must_not;
		}

		$http.post('/api/index/' + $scope.theindex.name + '/_search', requestBody).
		success(function(data) {
			$scope.results = data;
			for(var i in $scope.results.hits) {
				hit = $scope.results.hits[i];
				hit.roundedScore = $scope.roundScore(hit.score);
				hit.explanationString = $scope.expl(hit.explanation);
			}
			$scope.results.roundedTook = $scope.roundScore($scope.results.took);
		}).
		error(function(data, code) {

		});
	};
}