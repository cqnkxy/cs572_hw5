var app = angular.module('app', ['angularUtils.directives.dirPagination', 'ui.bootstrap', 'ngAnimate', 'ngSanitize']);

function searchViewController($scope, $http) {
	var shouldCorrect = false;
	var hasSearched = false;
	$scope.getSuggestions = function(val) {
		var words = val.split(/\s+/);
		var lastWord = words[words.length-1];
		shouldCorrect = false;
		hasSearched = false;
		return $http.get('suggest/', {
			params: {
				words: lastWord
			}
		}).then(function(response){
			var res = [];
			response.data.forEach(function(word){
				words[words.length-1] = word;
				res.push(words.join(" "));
			});
			return res;
		});
	};
	$scope.search = function() {
		$scope.results = [];
		$scope.query = $scope.query.split(/\s+/).join(" ");
		$http.get('search/', {
			params: {
				query: $scope.query,
				method: "default"
			}
		}).then(function(response){
			$scope.results = response.data;
			hasSearched = true;
		});
		$scope.corrected = "";
		$http.get('correct/', {
			params: {
				words: $scope.query,
			}
		}).then(function(response){
			if (response.data.toLowerCase() != $scope.query.toLowerCase()) {
				$scope.corrected = response.data;
				shouldCorrect = true;
			}
		});
	};
	$scope.correctSearch = function() {
		console.log(`query=${$scope.query}, corrected=${$scope.corrected}`);
		shouldCorrect = false;
		$scope.query = $scope.corrected;
		$scope.results = [];
		$http.get('search/', {
			params: {
				query: $scope.corrected,
				method: "default"
			}
		}).then(function(response){
			$scope.results = response.data;
		});	
	};
	$scope.checkCorrect = function() {
		return shouldCorrect;
	}
	$scope.emptyResults = function() {
		return hasSearched && $scope.results.length == 0;
	};
}

app.controller('searchViewController', searchViewController);
