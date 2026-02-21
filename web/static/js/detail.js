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
	// Build breadcrumb if this is a sub-container
	let breadcrumbHtml = '';
	if (tote.parent_id) {
		breadcrumbHtml = `
			<div style="margin-bottom: 1rem; padding: 0.5rem; background: #f0f0f0; border-radius: 5px;">
				<a href="/tote/${tote.parent_id}" style="color: #2196F3; text-decoration: none;">← Back to Parent Container</a>
			</div>
		`;
	}

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

	// Build sub-containers section (only for top-level containers)
	let childrenHtml = '';
	if (tote.depth === 0 && tote.children && tote.children.length > 0) {
		childrenHtml = `
			<div style="margin: 2rem 0;">
				<h3 style="margin-bottom: 1rem;">Sub-Containers (${tote.children.length})</h3>
				<div style="display: grid; grid-template-columns: repeat(auto-fill, minmax(250px, 1fr)); gap: 1rem;">
					${tote.children.map(child => {
						const childImage = child.images && child.images.length > 0 
							? `<img src="${child.images[0].image_data}" style="width: 100%; height: 120px; object-fit: cover; border-radius: 5px 5px 0 0;">`
							: '<div style="width: 100%; height: 120px; background: #f5f5f5; border-radius: 5px 5px 0 0; display: flex; align-items: center; justify-content: center;">📦</div>';
						
						return `
							<div onclick="window.location.href='/tote/${child.id}'" style="cursor: pointer; border: 1px solid #ddd; border-radius: 5px; overflow: hidden; transition: transform 0.2s;" onmouseover="this.style.transform='translateY(-2px)'" onmouseout="this.style.transform='translateY(0)'">
								${childImage}
								<div style="padding: 1rem;">
									<div style="font-weight: 600; margin-bottom: 0.5rem;">${child.name}</div>
									<div style="font-size: 0.85rem; color: #666;">${child.qr_code}</div>
									${child.description ? `<div style="font-size: 0.85rem; color: #888; margin-top: 0.3rem;">${child.description}</div>` : ''}
								</div>
							</div>
						`;
					}).join('')}
				</div>
			</div>
		`;
	}

	// Show "Add Sub-Container" button only for top-level containers
	const addSubContainerBtn = tote.depth === 0 
		? `<button class="btn btn-primary" onclick="window.location.href='/add?parent_id=${tote.id}'" style="margin-left: 0.5rem;">
				➕ Add Sub-Container
			</button>`
		: '';

	// Hide QR code section for sub-containers
	const qrSectionHtml = tote.depth === 0 ? `
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
	` : `
		<div class="detail-info">
			${descriptionHtml}
			${locationHtml}
			<div class="detail-row">
				<label>Type</label>
				<div class="value" style="color: #FF9800; font-weight: 600;">📦 Sub-Container</div>
			</div>
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
	`;

	const html = `
		${breadcrumbHtml}
		<div class="detail-header">
			<h2>${tote.name}</h2>
			<div class="tote-qr-code" style="font-size: 1.1rem; margin-top: 0.5rem;">${tote.qr_code}</div>
			${addSubContainerBtn}
		</div>

		${qrSectionHtml}

		${imagesHtml}

		${itemsHtml}

		${childrenHtml}
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
	// Only generate QR code if the qrcode div exists (not shown for sub-containers)
	const qrcodeDiv = document.getElementById('qrcode');
	if (qrcodeDiv) {
		new QRCode(qrcodeDiv, {
			text: qrText,
			width: 150,
			height: 150,
			colorDark: '#000000',
			colorLight: '#ffffff',
			correctLevel: QRCode.CorrectLevel.H
		});
	}
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


