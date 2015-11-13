'use strict';

angular.module('boomerApp')
    .controller('MainCtrl', function($scope, $http) {
        $scope.awesomeThings = [];
        $scope.endpoint1 = '';
        $scope.endpooint1Method = 'GET';
        $scope.endpoint2 = '';
        $scope.endpooint2Method = 'GET';

        $scope.battle = function() {

            var req1 = {
                requests: 100,
                concurrent_requests: 10,
                method: $scope.endpooint1Method,
                endpoint: $scope.endpoint1
            };
            var req2 = {
                requests: 100,
                concurrent_requests: 10,
                method: $scope.endpooint2Method,
                endpoint: $scope.endpoint2
            };

            $http({
                url: 'http://localhost:8081/bench',
                method: 'post',
                data: req1
            }).success(function(data) {
                $scope.response1 = {
                	endpoint: $scope.endpoint1,
                	data: data
                };
                $scope.endpoint1 = '';
                generateGraph(data.histogram, '.response1');
            });
            $http({
                url: 'http://localhost:8081/bench',
                method: 'post',
                data: req2
            }).success(function(data) {
                $scope.response2 = {
                	endpoint: $scope.endpoint2,
                	data: data
                };
                $scope.endpoint2 = '';
                generateGraph(data.histogram, '.response2');
            });
        };

        function generateGraph(hist, selector) {
        	jQuery(selector).html('');
            // Set the dimensions of the canvas / graph
            var margin = {
                    top: 20,
                    right: 20,
                    bottom: 120,
                    left: 40
                },
                width = 600 - margin.left - margin.right,
                height = 600 - margin.top - margin.bottom;
            // Parse the date / time
            var parseDate = d3.time.format("%L").parse;

            var x = d3.scale.ordinal().rangeRoundBands([0, width], .05);

            var y = d3.scale.linear().range([height, 0]);

            var xAxis = d3.svg.axis()
                .scale(x)
                .orient("bottom")
                .ticks(100);

            var yAxis = d3.svg.axis()
                .scale(y)
                .orient("left")
                .ticks(10);

            // Define the line
            var priceline = d3.svg.line()
                .x(function(d) {
                    return x(d.bucket);
                })
                .y(function(d) {
                    return y(d.count);
                });

            // Adds the svg canvas
            var svg = d3.select(selector)
                .attr('width', width + margin.left + margin.right)
                .attr('height', height + margin.top + margin.bottom)
                .append('g')
                .attr('transform',
                    'translate(' + margin.left + ',' + margin.top + ')');


            hist.forEach(function(b) {
                b.bucket = b.bucket.toFixed(2).toString();
                b.count = +b.count;
            });

            x.domain(hist.map(function(d) {
                return d.bucket;
            }));
            y.domain([0, d3.max(hist, function(d) {
                return d.count;
            })]);

            svg.append("g")
                .attr("class", "x axis")
                .attr("transform", "translate(0," + height + ")")
                .call(xAxis)
                .selectAll("text")
                .style("text-anchor", "end")
                .attr("dx", "-.8em")
                .attr("dy", "-.55em")
                .attr("transform", "rotate(-90)");
            svg.append("g")
                .attr("class", "y axis")
                .call(yAxis)
                .append("text")
                .attr("transform", "rotate(-90)")
                .attr("y", 6)
                .attr("dy", ".71em")
                .style("text-anchor", "end")
                .text("Request Count");
            svg.selectAll("bar")
                .data(hist)
                .enter().append("rect")
                .style("fill", "steelblue")
                .attr("x", function(d) {
                    return x(d.bucket);
                })
                .attr("width", x.rangeBand())
                .attr("y", function(d) {
                    return y(d.count);
                })
                .attr("height", function(d) {
                    return height - y(d.count);
                });
        }

        function getQueryVariables(query) {
            // var query = window.location.search.substring(1);
            var vars = query.split('&');
            var pairs = [];
            for (var i = 0; i < vars.length; i++) {
                var pair = vars[i].split('=');
                pairs.push({
                    key: decodeURIComponent(pair[0]),
                    value: decodeURIComponent(pair[1])
                });
            }
            return pairs;
        }

    });