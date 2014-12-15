var app = angular.module('snippetsApp', ['ngSanitize', 'ngCookies', 'ui.router', 'ui.bootstrap', 'ui.ace'])
    .run(['$rootScope', '$state', '$stateParams', '$cookies',
            function($rootScope, $state, $stateParams, $cookies) {
                // Make the state accessible from the root scope
                $rootScope.$state = $state;
                $rootScope.$stateParams = $stateParams;
                // Read User object from the cookie
                var c = B64.decode($cookies['snippets-auth']);
                c = c.split("|")[1];
                c = B64.decode(c);
                c = c.substring(c.indexOf("{"));
                $rootScope.user = JSON.parse(c);
            }
        ]);

app.config(['$stateProvider', '$urlRouterProvider',
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

            .state('snippets-new', {
                url: '/snippets/new',
                views: angular.extend({}, views, {'main': {
                    templateUrl: 'html/snippets.edit.html',
                    controller: 'SnippetEditCtrl'
                }})})

            .state('all-teams', {
                url: '/all-teams',
                views: angular.extend({}, views, {'main': {
                    templateUrl: 'html/all-teams.html'
                }})});
    }]);

app.controller('NavbarCtrl', function($scope, $cookies) {
    $scope.user = $scope.$root.user;
    window.user = $scope.user;
});

app.controller('MenuCtrl', function($scope, $location) {
    $scope.state = $scope.$root.$state.current.name;
});

app.controller('SnippetEditCtrl', function($scope) {

    $scope.preview = "";

    $scope.aceLoaded = function(editor) {
        editor.setTheme("ace/theme/textmate");
        editor.getSession().setMode("ace/mode/markdown");
        editor.focus();
        $scope.editor = editor;
    };

    $scope.renderPreview = function() {
        if (!$scope.editor) {
            return;
        }
        $scope.preview = markdown.toHTML($scope.editor.getValue());
        if (!$scope.preview) {
            $scope.preview = '<span class="italic">no preview available</span>';
        }
    }
});
