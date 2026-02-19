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
	// Build images gallery HTML
	let imagesHtml = '';
	if (tote.images && tote.images.length > 0) {
		imagesHtml = '<div class="images-gallery" style="display: grid; grid-template-columns: repeat(auto-fill, minmax(200px, 1fr)); gap: 15px; margin: 2rem 0;">';
		tote.images.forEach(img => {
			imagesHtml += `
				<div class="image-item" style="position: relative; cursor: pointer;" onclick="viewFullImage('${img.image_data}')">
					<img src="${img.image_data}" class="detail-image" alt="${tote.name}" style="width: 100%; height: 200px; object-fit: cover;">
				</div>
			`;
		});
		imagesHtml += '</div>';
	}

	const descriptionHtml = tote.description 
		? `<div class="detail-row">
				<label>Description</label>
				<div class="value">${tote.description}</div>
			</div>`
		: '';

	const locationHtml = tote.location 
		? `<div class="detail-row">
				<label>Location</label>
				<div class="value">${tote.location}</div>
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
				${locationHtml}
				<div class="detail-row">
					<label>Total Images</label>
					<div class="value">${tote.images ? tote.images.length : 0}</div>
				</div>
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

		${imagesHtml}

		${itemsHtml}
	`;

	document.getElementById('tote-detail').innerHTML = html;
	
	// Add modal for full-size image viewing
	if (!document.getElementById('image-modal')) {
		const modal = document.createElement('div');
		modal.id = 'image-modal';
		modal.style.display = 'none';
		modal.style.position = 'fixed';
		modal.style.top = '0';
		modal.style.left = '0';
		modal.style.width = '100%';
		modal.style.height = '100%';
		modal.style.backgroundColor = 'rgba(0, 0, 0, 0.9)';
		modal.style.zIndex = '10000';
		modal.style.cursor = 'pointer';
		modal.innerHTML = `
			<div style="position: relative; width: 100%; height: 100%; display: flex; align-items: center; justify-content: center;">
				<span style="position: absolute; top: 20px; right: 35px; color: #f1f1f1; font-size: 40px; font-weight: bold; cursor: pointer;" onclick="closeImageModal()">&times;</span>
				<img id="modal-image" style="max-width: 95%; max-height: 95%; object-fit: contain;">
			</div>
		`;
		document.body.appendChild(modal);
		
		// Close modal when clicking outside image
		modal.addEventListener('click', closeImageModal);
	}
}

function viewFullImage(imageSrc) {
	const modal = document.getElementById('image-modal');
	const modalImg = document.getElementById('modal-image');
	modal.style.display = 'block';
	modalImg.src = imageSrc;
}

function closeImageModal() {
	document.getElementById('image-modal').style.display = 'none';
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


