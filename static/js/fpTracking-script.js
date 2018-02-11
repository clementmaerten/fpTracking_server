

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

			updateProgressMessage(data.CurrentVisitFrequency);

			updateProgressBar(data.Progression);

			updateAllGraphics(data);

			if (data.Progression >= 100) {
				clearInterval(checkIntervalId);
				stopProgressBar();
				deleteProgressMessage();

				displayExecutingTime(data.ExecutingTime);
				displayBackButton();
			}
		},
		error: () => {
			//alert("Error in checkProgression");
		}
	});
}

//BACK BUTTON
function displayBackButton() {
	$('#backButtonDivId').show();
}

$('#backButtonDivId').click(() => {
	resetProgressBar();
	resetGraphics();
	deleteExecutingTime();

	$('#backButtonDivId').hide();
	$('#fpTrackingParallelResultsId').hide();

	$('#fpTrackingParallelForm').show();
});

//PROGRESS MESSAGE
function updateProgressMessage(currentVisitFrequency) {
	$('#progressMessageId').show();
	$('#progressMessageId').html('Computing visit frequency '+currentVisitFrequency+' ...');
}

function deleteProgressMessage() {
	$('#progressMessageId').hide();
	$('#progressMessageId').html('');
}

//EXECUTING TIME
function displayExecutingTime(executingTime) {
	$('#ExecutingTimeId').show();
	if (executingTime <= 1) {
		$('#ExecutingTimeId').html('Executed in '+executingTime+' second');
	} else {
		$('#ExecutingTimeId').html('Executed in '+executingTime+' seconds');
	}
}

function deleteExecutingTime() {
	$('#ExecutingTimeId').hide();
	$('#ExecutingTimeId').html('');
}

//PROGRESS BAR
function updateProgressBar(progression) {
	$('#progressBarId').attr('aria-valuenow',progression).css('width',progression+'%').html(progression+'%');
}

function stopProgressBar() {
	$('#progressBarId').removeClass('progress-bar-animated').removeClass('progress-bar-striped');
}

function resetProgressBar() {
	$('#progressBarId').attr('aria-valuenow','0').css('width','0%').html('0%');
	$('#progressBarId').addClass('progress-bar-animated').addClass('progress-bar-striped');
}


//GRAPHICS
var rawAndRawMaxDaysFrequencyGraph;
var nbIdsFrequencyGraph;
var ownershipFrequencyGraph;

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
	    	showInLegend: false,
	        name: 'Average',
	        data: []
	    }, {
	    	showInLegend: false,
	    	name: 'Maximum average',
	    	data: []
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
	    	showInLegend: false,
	        name: 'Rule-based',
	        data: []
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
	    	showInLegend: false,
	        name: 'Rule-based',
	        data: []
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

function updateAllGraphics(data) {
	for (let i = 0; i < data.AverageTrackingTimeGraph.length; i++) {
		rawAndRawMaxDaysFrequencyGraph.series[0].addPoint([
			data.AverageTrackingTimeGraph[i].VisitFrequency,
			data.AverageTrackingTimeGraph[i].Value
		]);
		rawAndRawMaxDaysFrequencyGraph.series[1].addPoint([
			data.MaximumAverageTrackingTimeGraph[i].VisitFrequency,
			data.MaximumAverageTrackingTimeGraph[i].Value
		]);
		nbIdsFrequencyGraph.series[0].addPoint([
			data.NbIdsFrequencyGraph[i].VisitFrequency,
			data.NbIdsFrequencyGraph[i].Value
		]);
		ownershipFrequencyGraph.series[0].addPoint([
			data.OwnershipFrequencyGraph[i].VisitFrequency,
			data.OwnershipFrequencyGraph[i].Value
		]);
	}
}

function resetGraphics() {
	rawAndRawMaxDaysFrequencyGraph.series[0].data = [];
	rawAndRawMaxDaysFrequencyGraph.series[1].data = [];

	nbIdsFrequencyGraph.series[0].data = [];

	ownershipFrequencyGraph.series[0].data = [];
}



//WHEN THE USER LEAVES THE PAGE
var hasUserLeftPage = false;

$(window).on('beforeunload', () => {
	//this will work only for Chrome
	sendStopRequest();
});

$(window).on('unload', () => {
	//this will work for other browsers
	sendStopRequest();
});

function sendStopRequest() {
	if (!hasUserLeftPage) {
		$.ajax({
			url: 'stop-tracking',
			type: 'POST',
			async: false,
			success: (data) => {
				hasUserLeftPage = true;
			},
			error: () => {
				//alert("Error in checkProgression");
			}
		});
	}
}