

$('#fpTrackingParallelForm').submit(() => {
	//verify parameters 
	//useful ?
	
	//send the request
	$.ajax({
		url: 'tracking_parallel/',
		type: 'POST',
		data: $('#fpTrackingParallelForm').serialize(),
		success: (data) => {
			//alert(data); 
		},
		error: (e) => {
			console.log(e);
			alert('La requête n\'a pas abouti'); 
		}
	});
	return false;
});