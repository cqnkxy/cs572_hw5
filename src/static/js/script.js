var app = angular.module('app', ['angularUtils.directives.dirPagination', 'ui.bootstrap', 'ngAnimate', 'ngSanitize']);

function searchViewController($scope, $http) {
	$scope.getSuggestions = function(val) {
		var words = val.split(/\s+/);
		var lastWord = words[words.length-1];
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
		$scope.results = []
		$http.get('search/', {
			params: {
				query: $scope.query,
				method: "default"
			}
		}).then(function(response){
			$scope.results = response.data;
		});
	};
}

app.controller('searchViewController', searchViewController);
