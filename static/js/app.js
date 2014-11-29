var postsApp = angular.module('postsApp', ['ngRoute']);

postsApp.config(['$routeProvider',
    function($routeProvider) {
        $routeProvider.
          when('/', {
            templateUrl: 'html/index.html',
            controller: 'IndexCtrl'
          }).
          otherwise({redirectTo: '/'});
    }]);

postsApp.controller('NavbarCtrl', function($scope, $http) {
    $http.get('/api/user').success(function(data) {
        $scope.user = data;
    });
});

postsApp.controller('IndexCtrl', function($scope) {
});
