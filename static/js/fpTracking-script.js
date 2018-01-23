

$('#fpTrackingParallelForm').submit(() => {
	$.ajax({
		url: 'testPost/',
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