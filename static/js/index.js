angular.module("webtail").controller("mainController", mainController);

mainController.$inject = ["$rootScope", "$scope", "$mdSidenav", "$http"]

function mainController($rootScope, $scope, $mdSidenav, $http) {
  var vm = this;

  vm.toggleSideNav = function toggleSideNav() {
    $mdSidenav('left').toggle()
  }

  vm.init = function init() {
    console.log("In the main controller")
    $scope.showCard=true;
    $http.get('user')
      .then(function(result){
        $rootScope.username = result.data["username"]
        $rootScope.isLoggedIn = result.data["isLoggedIn"]
        console.log("is logged in :", result.data)
      }, function(result) {
        console.log("Failed to get the username")
      })
  }

  vm.fontSize = ["10px", "11px", "12px", "14px", "16px", "18px", "20px", "22px", "24px"]
  $scope.currSize = vm.fontSize[2];

  $scope.open_connection = function(file) {
    console.log(file)
    $scope.showCard=false;
    // $scope.$apply()
    angular.element(document.querySelector("#filename")).html("File: " + file)
    var container = angular.element(document.querySelector("#container"))
    var ws;
    if (window.WebSocket === undefined) {
        container.append("Your browser does not support WebSockets");
        return;
    } else {
        ws = initWS(file);
    }
    vm.toggleSideNav()
  }

  function initWS(file) {
    var socket = new WebSocket("ws://"+window.location.hostname+":" + window.location.port + "/ws/" + btoa(file));
    var container = angular.element(document.querySelector("#container"));

    // clear the contents
    container.html("")
    socket.onopen = function() {
        container.append("<p><b>Tailing file: " + file + "</b></p>");
    };
    socket.onmessage = function (e) {
        container.append(e.data.trim()+"<br>");
    }
    socket.onclose = function () {
        container.append("<p>Connection Closed to WebSocket, tail stopped</p>");
    }
    return socket;
  }

  $scope.logout = function() {
    for (i = 0; i < document.forms.length; i++) {
      if (document.forms[i].id == "logoutForm") {
        document.forms[i].submit()
        return;
      }
    }
  }

  vm.init();
}
