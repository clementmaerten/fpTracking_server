

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

$('#fpTrackingParallelForm').submit(() => {
	
	//verify parameters
	const number = parseInt($('#fpTrackingParallelNumberId').val());
	const minNbPerUser = parseInt($('#fpTrackingParallelMinNbPerUserId').val());
	const goroutineNumber = parseInt($('#fpTrackingParallelGoroutineNumberId').val());

	if (isNaN(number) || isNaN(minNbPerUser)  || isNaN(goroutineNumber)){
		triggerErrorMessage('fpTrackingParallelResultsErrorAlertId','Invalid format for parameters');
	} else if (number <= 0 || minNbPerUser <=0 || goroutineNumber <=0) {
		triggerErrorMessage('fpTrackingParallelResultsErrorAlertId','nul or negative parameters');
	} else {

		//Begin the check of progression every second
		const checkIntervalId = setInterval(checkProgression,1000);

		//send the request
		$.ajax({
			url: 'tracking-parallel',
			type: 'POST',
			data: $('#fpTrackingParallelForm').serialize(),
			success: (data) => {
				//stop checkProgression
				clearInterval(checkIntervalId);

				//launch for the last time the checkProgression function
				checkProgression();

				//clear the error alerts
				deleteErrorMessage('fpTrackingParallelResultsErrorAlertId');
			},
			error: () => {
				//stop checkProgression
				clearInterval(checkIntervalId);

				triggerErrorMessage('fpTrackingParallelResultsErrorAlertId','The server wasn\'t able to process the request');
			}
		});
	}

	//We return false so that the function doesn't refresh the page
	return false;
});

function checkProgression (){
	$.ajax({
		url: 'check-progression',
		type: 'POST',
		success: (data) => {

		},
		error: () => {
			alert("Error in checkProgression");
		}
	});
}