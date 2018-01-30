

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
	let number = parseInt($('#fpTrackingParallelNumberId').val());
	let minNbPerUser = parseInt($('#fpTrackingParallelMinNbPerUserId').val());
	let goroutineNumber = parseInt($('#fpTrackingParallelGoroutineNumberId').val());

	if (isNaN(number) || isNaN(minNbPerUser)  || isNaN(goroutineNumber)){
		triggerErrorMessage('fpTrackingParallelResultsErrorAlertId','Invalid format for parameters');
	} else if (number <= 0 || minNbPerUser <=0 || goroutineNumber <=0) {
		triggerErrorMessage('fpTrackingParallelResultsErrorAlertId','nul or negative parameters');
	} else {
		//send the request
		$.ajax({
			url: 'tracking-parallel',
			type: 'POST',
			data: $('#fpTrackingParallelForm').serialize(),
			success: (data) => {
				//clear the error alerts
				deleteErrorMessage('fpTrackingParallelResultsErrorAlertId');
			},
			error: (e) => {
				triggerErrorMessage('fpTrackingParallelResultsErrorAlertId','The server wasn\'t able to process the request');
			}
		});
	}

	//We return false so that the function doesn't refresh the page
	return false;
});