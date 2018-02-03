

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

function updateProgressBar(progression) {
	$('#progressBarId').attr('aria-valuenow',progression).css('width',progression+'%').html(progression+'%');
}

function stopProgressBar() {
	$('#progressBarId').removeClass('progress-bar-animated').removeClass('progress-bar-striped');
}