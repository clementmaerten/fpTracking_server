

$('#fpTrackingParallelForm').submit(() => {

	//clear the display
	$('#fpTrackingParallelResultsErrorAlertId').hide();
	
	//verify parameters
	let number = parseInt($('#fpTrackingParallelNumberId').val());
	let minNbPerUser = parseInt($('#fpTrackingParallelMinNbPerUserId').val());
	if (isNaN(number) || isNaN(minNbPerUser)){
		alert("Invalid format for parameters")
	} else if (number <= 0 || minNbPerUser <=0) {
		alert("nul or negative parameters");
	} else {
		//send the request
		$.ajax({
			url: 'tracking_parallel/',
			type: 'POST',
			data: $('#fpTrackingParallelForm').serialize(),
			success: (data) => {
				//alert(data); 
			},
			error: (e) => {
				$('#fpTrackingParallelResultsErrorAlertId').show();
				$('#fpTrackingParallelResultsErrorAlertId').html('The server wasn\'t able to process the request');
			}
		});
	}

	return false;
});