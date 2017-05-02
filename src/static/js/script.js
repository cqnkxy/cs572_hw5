var app = angular.module('app', [
	'angularUtils.directives.dirPagination', 
	'ui.bootstrap', 'ngAnimate', 'ngSanitize'
]);

app.controller("searchViewController", function($scope, $http) {
	var hasCorrected = false;
	var hasSearched = false;
	$scope.getSuggestions = function(val) {
		var words = val.split(/\s+/);
		var lastWord = words[words.length-1];
		hasCorrected = false;
		hasSearched = false;
		return $http.get('suggest/', {
			params: {
				words: lastWord
			}
		}).then(function(response){
			var res = [];
			response.data.forEach(function(word){
				// if (word.startsWith(lastWord)){
				// 	words[words.length-1] = word;
				// 	res.push(words.join(" "));
				// }
				words[words.length-1] = word;
				res.push(words.join(" "));
			});
			return res;
		});
	};
	$scope.search = function() {
		$scope.results = [];
		$scope.query = $scope.query.split(/\s+/).join(" ");
		$scope.corrected = "";
		$http.get('correct/', {
			params: {
				words: $scope.query,
			}
		}).then(function(response){
			$scope.corrected = response.data;
			if (response.data.toLowerCase() != $scope.query.toLowerCase()) {
				hasCorrected = true;
			}
			$http.get('search/', {
				params: {
					query: $scope.corrected,
					method: "default"
				}
			}).then(function(response){
				$scope.results = response.data;
				hasSearched = true;
			});
		});
		
	};
	$scope.insteadSearch = function() {
		console.log(`query=${$scope.query}, corrected=${$scope.corrected}`);
		hasCorrected = false;
		$scope.results = [];
		$http.get('search/', {
			params: {
				query: $scope.query,
				method: "default"
			}
		}).then(function(response){
			$scope.results = response.data;
		});
	};
	$scope.checkCorrect = function() {
		return hasCorrected;
	}
	$scope.emptyResults = function() {
		return hasSearched && $scope.results.length == 0;
	};
})
.filter('highlight', function($sce) {
    return function(text, phrase) {
      if (phrase) text = text.replace(new RegExp('('+phrase.split(" ").join("|")+')', 'gi'),
        '<span class="highlighted">$1</span>')

      return $sce.trustAsHtml(text)
    }
})
