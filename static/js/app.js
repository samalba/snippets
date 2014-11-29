var postsApp = angular.module('postsApp', ['ui.router'])
    .run(['$rootScope', '$state', '$stateParams',
            function($rootScope, $state, $stateParams) {
                // Make the state accessible from the root scope
                $rootScope.$state = $state;
                $rootScope.$stateParams = $stateParams;
            }
        ]);

postsApp.config(['$stateProvider', '$urlRouterProvider',
    function($stateProvider, $urlRouterProvider) {

        $urlRouterProvider
            .otherwise('/');

        var views = {
            'navbar': {
                templateUrl: 'html/navbar.html',
                controller: 'NavbarCtrl'
            },
            'menu': {
                templateUrl: 'html/menu.html',
                controller: 'MenuCtrl'
            }
        };
        $stateProvider
            .state('index', {
                url: '/',
                views: angular.extend({}, views, {'main': {
                    templateUrl: 'html/index.html'
                }})})
            .state('all-teams', {
                url: '/all-teams',
                views: angular.extend({}, views, {'main': {
                    templateUrl: 'html/all-teams.html'
                }})});
    }]);

postsApp.controller('NavbarCtrl', function($scope, $http) {
    $http.get('/api/user').success(function(data) {
        $scope.user = data;
    });
});

postsApp.controller('MenuCtrl', function($scope, $location) {
    $scope.state = $scope.$root.$state.current.name;
    window.d = $scope.$root;
});
