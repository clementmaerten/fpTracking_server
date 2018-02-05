

function triggerErrorMessage(errorMessageId, message) {
	let objectId = '#'+errorMessageId;
	$(objectId).html(message);
	$(objectId).show();
}

function deleteErrorMessage(errorMessageId) {
	let objectId = '#'+errorMessageId;
	$(objectId).hide();
	$(objectId).html('');
}

var checkIntervalId;

$('#fpTrackingParallelForm').submit(() => {
	
	//verify parameters
	const number = parseInt($('#fpTrackingParallelNumberId').val());
	const minNbPerUser = parseInt($('#fpTrackingParallelMinNbPerUserId').val());
	const goroutineNumber = parseInt($('#fpTrackingParallelGoroutineNumberId').val());

	if (isNaN(number) || isNaN(minNbPerUser)  || isNaN(goroutineNumber)){
		triggerErrorMessage('fpTrackingParallelResultsErrorAlertId','Invalid format for parameters');
	} else if (number <= 0 || minNbPerUser <=0 || goroutineNumber <=0) {
		triggerErrorMessage('fpTrackingParallelResultsErrorAlertId','Nul or negative parameters');
	} else {

		//send the request
		$.ajax({
			url: 'tracking-parallel',
			type: 'POST',
			data: $('#fpTrackingParallelForm').serialize(),
			success: (data) => {
				//clear the error alerts
				deleteErrorMessage('fpTrackingParallelResultsErrorAlertId');

				//hide the form
				$('#fpTrackingParallelForm').hide();

				//show the results div
				$('#fpTrackingParallelResultsId').show();

				//Begin the check of progression every 5 seconds
				checkIntervalId = setInterval(checkProgression,5000);

				//Display the graphics
				displayGraphics();
			},
			error: () => {
				triggerErrorMessage('fpTrackingParallelResultsErrorAlertId','The server wasn\'t able to process the request');
			}
		});
	}

	//We return false so that the function doesn't refresh the page
	return false;
});

function checkProgression() {
	$.ajax({
		url: 'check-progression',
		type: 'POST',
		success: (data) => {

			updateProgressBar(data.Progression);

			//updateAllGraphics();

			if (data.Progression >= 100) {
				clearInterval(checkIntervalId);
				stopProgressBar();
			}
		},
		error: () => {
			//alert("Error in checkProgression");
		}
	});
}


//PROGRESS BAR
function updateProgressBar(progression) {
	$('#progressBarId').attr('aria-valuenow',progression).css('width',progression+'%').html(progression+'%');
}

function stopProgressBar() {
	$('#progressBarId').removeClass('progress-bar-animated').removeClass('progress-bar-striped');
}


//GRAPHICS
var rawAndRawMaxDaysFrequencyGraph;
var nbIdsFrequencyGraph;
var ownershipFrequencyGraph;
var myData = [43934, 52503, 57177, 69658, 97031, 119931, 137133, 154175];
var myData2 = [154175, 137133, 119931, 97031, 69658, 57177, 52503, 43934];

function displayGraphics() {
	rawAndRawMaxDaysFrequencyGraph = Highcharts.chart('rawAndRawMaxDaysFrequencyGraphId', {

	    title: {
	        text: 'Days Frequency graph'
	    },

	    yAxis: {
	        title: {
	            text: 'Average tracking time (days)'
	        }
	    },

	    xAxis: {
	        title: {
	            text: 'Collect frequency (days)'
	        }
	    },

	    plotOptions: {
	        series: {
	            label: {
	                connectorAllowed: false
	            },
	            pointStart: 0
	        }
	    },

	    series: [{
	        name: 'Average',
	        data: myData
	    }, {
	    	name: 'Maximum average',
	    	data: myData2
	    }],

	    responsive: {
	        rules: [{
	            condition: {
	                maxWidth: 500
	            },
	            chartOptions: {
	                legend: {
	                    layout: 'horizontal',
	                    align: 'center',
	                    verticalAlign: 'bottom'
	                }
	            }
	        }]
	    }

	});

	nbIdsFrequencyGraph = Highcharts.chart('nbIdsFrequencyGraphId', {

	    title: {
	        text: 'Number of ids Frequency graph'
	    },

	    yAxis: {
	        title: {
	            text: 'Number of ids per user'
	        }
	    },

	    xAxis: {
	        title: {
	            text: 'Collect frequency (days)'
	        }
	    },

	    plotOptions: {
	        series: {
	            label: {
	                connectorAllowed: false
	            },
	            pointStart: 0
	        }
	    },

	    series: [{
	        name: 'Rule-based',
	        data: myData
	    }],

	    responsive: {
	        rules: [{
	            condition: {
	                maxWidth: 500
	            },
	            chartOptions: {
	                legend: {
	                    layout: 'horizontal',
	                    align: 'center',
	                    verticalAlign: 'bottom'
	                }
	            }
	        }]
	    }

	});

	ownershipFrequencyGraph = Highcharts.chart('ownershipFrequencyGraphId', {

	    title: {
	        text: 'Ownership Frequency graph'
	    },

	    yAxis: {
	        title: {
	            text: 'Average ownership'
	        }
	    },

	    xAxis: {
	        title: {
	            text: 'Collect frequency (days)'
	        }
	    },

	    plotOptions: {
	        series: {
	            label: {
	                connectorAllowed: false
	            },
	            pointStart: 0
	        }
	    },

	    series: [{
	        name: 'Rule-based',
	        data: myData
	    }],

	    responsive: {
	        rules: [{
	            condition: {
	                maxWidth: 500
	            },
	            chartOptions: {
	                legend: {
	                    layout: 'horizontal',
	                    align: 'center',
	                    verticalAlign: 'bottom'
	                }
	            }
	        }]
	    }

	});
}

function updateAllGraphics() {
	myData[1] += 80000;
	myData[5] += 10000;
	//myData[8] += 10000;
	rawDaysFrequencyGraph.series[0].setData(myData,false);
	//rawDaysFrequencyGraph.series[0].data[5].update(test);
}