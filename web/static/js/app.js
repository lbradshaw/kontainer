// app.js - Main dashboard functionality

let allTotes = [];

// Load totes on page load
document.addEventListener('DOMContentLoaded', function() {
	loadTotes();
	setupSearch();
});

function loadTotes() {
	fetch('/api/totes')
		.then(response => response.json())
		.then(totes => {
			allTotes = totes || [];
			updateStats();
			displayTotes(allTotes);
		})
		.catch(error => {
			console.error('Error loading totes:', error);
			document.getElementById('totes-grid').innerHTML = 
				'<div class="loading">Error loading totes</div>';
		});
}

function updateStats() {
	document.getElementById('total-totes').textContent = allTotes.length;
}

function displayTotes(totes) {
	const grid = document.getElementById('totes-grid');
	const emptyState = document.getElementById('empty-state');

	if (!totes || totes.length === 0) {
		grid.style.display = 'none';
		emptyState.style.display = 'block';
		return;
	}

	grid.style.display = 'grid';
	emptyState.style.display = 'none';

	grid.innerHTML = totes.map(tote => {
		// Use first image from images array, fallback to legacy image_path
		let imageHtml = '';
		if (tote.images && tote.images.length > 0) {
			const imageCount = tote.images.length > 1 ? `<span class="image-count">+${tote.images.length - 1}</span>` : '';
			imageHtml = `
				<div class="tote-image-container" onmouseenter="showImagePreview(event, ${tote.id})" onmouseleave="hideImagePreview()">
					<img src="${tote.images[0].image_data}" class="tote-image" alt="${tote.name}">
					${imageCount}
				</div>`;
		} else if (tote.image_path) {
			imageHtml = `<img src="${tote.image_path}" class="tote-image" alt="${tote.name}">`;
		}
		
		const descriptionHtml = tote.description
			? `<div class="tote-description">${tote.description}</div>`
			: '';

		const locationHtml = tote.location
			? `<div class="tote-location" style="font-size: 0.85rem; color: #666; margin: 0.3rem 0;">📍 ${tote.location}</div>`
			: '';

		const itemsPreview = tote.items 
			? `<div class="tote-items-preview">${tote.items.split('\n').slice(0, 3).join('\n')}</div>`
			: '';

		return `
			<div class="tote-card" onclick="window.location.href='/tote/${tote.id}'">
				<div class="tote-card-header">
					<div>
						<h3>${tote.name}</h3>
						<div class="tote-qr-code">${tote.qr_code}</div>
					</div>
				</div>
				${imageHtml}
				${descriptionHtml}
				${locationHtml}
				${itemsPreview}
			</div>
		`;
	}).join('');
}

let imagePreviewTimeout;

function showImagePreview(event, toteId) {
	event.stopPropagation();
	
	// Find the tote
	const tote = allTotes.find(t => t.id === toteId);
	if (!tote || !tote.images || tote.images.length === 0) return;
	
	// Clear any existing timeout
	if (imagePreviewTimeout) {
		clearTimeout(imagePreviewTimeout);
		imagePreviewTimeout = null;
	}
	
	// Remove existing preview
	const existingPreview = document.getElementById('image-preview-modal');
	if (existingPreview) {
		existingPreview.remove();
	}
	
	// Create modal
	const modal = document.createElement('div');
	modal.id = 'image-preview-modal';
	modal.className = 'image-modal';
	
	const imagesHtml = tote.images.map((img, idx) => 
		`<img src="${img.image_data}" alt="${tote.name}" onclick="enlargeImage(event, ${toteId}, ${idx})">`
	).join('');
	
	modal.innerHTML = `
		<div class="image-modal-content">
			<div class="image-modal-header">
				<span>${tote.name} - ${tote.images.length} Image(s)</span>
			</div>
			<div class="image-modal-grid">${imagesHtml}</div>
		</div>
	`;
	
	// Keep modal visible when hovering over it
	modal.addEventListener('mouseenter', function() {
		if (imagePreviewTimeout) {
			clearTimeout(imagePreviewTimeout);
			imagePreviewTimeout = null;
		}
	});
	
	modal.addEventListener('mouseleave', function() {
		hideImagePreview();
	});
	
	document.body.appendChild(modal);
}

function hideImagePreview() {
	if (imagePreviewTimeout) {
		clearTimeout(imagePreviewTimeout);
	}
	imagePreviewTimeout = setTimeout(() => {
		const modal = document.getElementById('image-preview-modal');
		if (modal) {
			modal.remove();
		}
	}, 300);
}

function enlargeImage(event, toteId, imageIndex) {
	event.stopPropagation();
	
	const tote = allTotes.find(t => t.id === toteId);
	if (!tote || !tote.images || !tote.images[imageIndex]) return;
	
	// Create enlarged image modal
	const enlargedModal = document.createElement('div');
	enlargedModal.className = 'image-enlarged-modal';
	enlargedModal.innerHTML = `
		<div class="image-enlarged-content">
			<button class="image-close-btn" onclick="this.closest('.image-enlarged-modal').remove()">✕</button>
			<img src="${tote.images[imageIndex].image_data}" alt="${tote.name}">
		</div>
	`;
	
	// Close on background click
	enlargedModal.addEventListener('click', function(e) {
		if (e.target === enlargedModal) {
			enlargedModal.remove();
		}
	});
	
	document.body.appendChild(enlargedModal);
}

function setupSearch() {
	const searchInput = document.getElementById('search');
	searchInput.addEventListener('input', function(e) {
		const query = e.target.value.toLowerCase();
		
		if (!query) {
			displayTotes(allTotes);
			return;
		}

		const filtered = allTotes.filter(tote => {
			return tote.name.toLowerCase().includes(query) ||
				(tote.description && tote.description.toLowerCase().includes(query)) ||
				(tote.items && tote.items.toLowerCase().includes(query)) ||
				tote.qr_code.toLowerCase().includes(query);
		});

		displayTotes(filtered);
	});
}

// Export data as JSON
function exportData() {
	window.location.href = '/api/export';
}

// Import data from JSON file
async function importData(event) {
	const file = event.target.files[0];
	if (!file) return;

	// Confirm import
	if (!confirm(`Import data from ${file.name}?\n\nThis will ADD the totes from the file to your existing inventory.`)) {
		event.target.value = ''; // Reset file input
		return;
	}

	try {
		const text = await file.text();
		const totes = JSON.parse(text);

		const response = await fetch('/api/import', {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: text
		});

		if (!response.ok) {
			throw new Error('Import failed');
		}

		const result = await response.json();
		alert(`Successfully imported ${result.imported} tote(s)!`);
		
		// Reload the page to show new totes
		window.location.reload();
	} catch (error) {
		console.error('Error importing data:', error);
		alert('Error importing data. Please check the file format.');
	} finally {
		event.target.value = ''; // Reset file input
	}
}
