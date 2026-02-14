// detail.js - Tote detail page functionality

document.addEventListener('DOMContentLoaded', function() {
	loadToteDetail();
});

function loadToteDetail() {
	const pathParts = window.location.pathname.split('/');
	const toteId = pathParts[pathParts.length - 1];

	fetch(`/api/tote/${toteId}`)
		.then(response => {
			if (!response.ok) {
				throw new Error('Tote not found');
			}
			return response.json();
		})
		.then(tote => {
			displayToteDetail(tote);
			generateQRCode(tote.qr_code);
		})
		.catch(error => {
			console.error('Error loading tote:', error);
			document.getElementById('tote-detail').innerHTML = 
				'<div class="loading">Error loading tote details</div>';
		});
}

function displayToteDetail(tote) {
	const imageHtml = tote.image_path 
		? `<img src="${tote.image_path}" class="detail-image" alt="${tote.name}">`
		: '';

	const descriptionHtml = tote.description 
		? `<div class="detail-row">
				<label>Description</label>
				<div class="value">${tote.description}</div>
			</div>`
		: '';

	const itemsHtml = tote.items 
		? `<div class="detail-row">
				<label>Items</label>
				<div class="items-list">${tote.items}</div>
			</div>`
		: '';

	const html = `
		<div class="detail-header">
			<h2>${tote.name}</h2>
			<div class="tote-qr-code" style="font-size: 1.1rem; margin-top: 0.5rem;">${tote.qr_code}</div>
		</div>

		<div class="detail-qr-section">
			<div class="detail-qr-code">
				<div id="qrcode"></div>
				<button class="btn btn-primary" onclick="window.location.href='/print-label/${tote.id}'" style="margin-top: 1rem;">
					🖨️ Print Label
				</button>
			</div>
			<div class="detail-info">
				${descriptionHtml}
				<div class="detail-row">
					<label>Created</label>
					<div class="value">${new Date(tote.created_at).toLocaleDateString()}</div>
				</div>
				<div class="detail-row">
					<label>Last Updated</label>
					<div class="value">${new Date(tote.updated_at).toLocaleDateString()}</div>
				</div>
			</div>
		</div>

		${imageHtml}

		${itemsHtml}
	`;

	document.getElementById('tote-detail').innerHTML = html;
}

function generateQRCode(qrText) {
	new QRCode(document.getElementById('qrcode'), {
		text: qrText,
		width: 150,
		height: 150,
		colorDark: '#000000',
		colorLight: '#ffffff',
		correctLevel: QRCode.CorrectLevel.H
	});
}

function deleteTote() {
	if (!confirm('Are you sure you want to delete this tote?')) {
		return;
	}

	const pathParts = window.location.pathname.split('/');
	const toteId = pathParts[pathParts.length - 1];

	fetch(`/api/tote/${toteId}`, {
		method: 'DELETE'
	})
	.then(response => {
		if (!response.ok) {
			throw new Error('Failed to delete tote');
		}
		window.location.href = '/';
	})
	.catch(error => {
		console.error('Error deleting tote:', error);
		alert('Error deleting tote');
	});
}
