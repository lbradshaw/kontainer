// form.js - Add/Edit tote form functionality

const isEditMode = window.location.pathname.includes('/edit');
let currentImagePaths = [];
let uploadedImages = []; // Changed to store base64 data and types

document.addEventListener('DOMContentLoaded', function() {
	setupForm();
	setupImagePreview();

	if (isEditMode) {
		loadToteData();
	}
});

function setupForm() {
	const form = document.getElementById('tote-form');
	form.addEventListener('submit', handleSubmit);
}

function setupImagePreview() {
	const imageInput = document.getElementById('image');
	imageInput.addEventListener('change', async function(e) {
		const files = e.target.files;
		if (files.length > 0) {
			const previewContainer = document.getElementById('image-preview');
			previewContainer.innerHTML = '<h4>Selected Images:</h4>';
			previewContainer.style.display = 'block';

			// Convert all images to base64
			uploadedImages = [];
			for (let i = 0; i < files.length; i++) {
				const file = files[i];
				
				// Convert to base64
				const reader = new FileReader();
				reader.onload = function(e) {
					const base64Data = e.target.result; // Already in data URI format
					
					uploadedImages.push({
						data: base64Data,
						type: file.type
					});
					
					// Show preview
					const img = document.createElement('img');
					img.src = base64Data;
					img.style.maxWidth = '150px';
					img.style.maxHeight = '150px';
					img.style.margin = '5px';
					img.style.border = '1px solid #ddd';
					img.style.borderRadius = '4px';
					previewContainer.appendChild(img);
				};
				reader.readAsDataURL(file);
			}
		}
	});
}

function loadToteData() {
	const toteId = document.getElementById('tote-id').value;
	
	fetch(`/api/tote/${toteId}`)
		.then(response => response.json())
		.then(tote => {
			document.getElementById('name').value = tote.name || '';
			document.getElementById('description').value = tote.description || '';
			document.getElementById('items').value = tote.items || '';
			document.getElementById('location').value = tote.location || '';
			
			// Store existing images
			currentImagePaths = tote.images || [];

			if (tote.images && tote.images.length > 0) {
				let imagesHtml = `<p><strong>Current images (${tote.images.length}):</strong></p><div style="display: flex; flex-wrap: wrap; gap: 10px;">`;
				tote.images.forEach(img => {
					imagesHtml += `
						<div style="position: relative; display: inline-block;">
							<img src="${img.image_data}" style="max-width: 150px; max-height: 150px; border: 1px solid #ddd; border-radius: 4px; display: block;">
							<button type="button" onclick="deleteImage(${img.id})" class="btn btn-danger" style="position: absolute; top: 5px; right: 5px; padding: 5px 10px; font-size: 0.8rem;">
								🗑️
							</button>
						</div>
					`;
				});
				imagesHtml += '</div>';
				document.getElementById('current-image').innerHTML = imagesHtml;
			}
		})
		.catch(error => {
			console.error('Error loading tote:', error);
			alert('Error loading tote data');
		});
}

async function handleSubmit(e) {
	e.preventDefault();

	const name = document.getElementById('name').value;
	const description = document.getElementById('description').value;
	const items = document.getElementById('items').value;
	const location = document.getElementById('location').value;

	if (isEditMode) {
		// For edit mode, just update the tote details
		const toteData = {
			name,
			description,
			items,
			location
		};

		const toteId = document.getElementById('tote-id').value;

		try {
			const response = await fetch(`/api/tote/${toteId}`, {
				method: 'PUT',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify(toteData)
			});

			if (!response.ok) {
				throw new Error('Failed to update tote');
			}

			// Add new images if any were uploaded (send base64 data)
			for (const img of uploadedImages) {
				await fetch(`/api/tote/${toteId}/add-image`, {
					method: 'POST',
					headers: {
						'Content-Type': 'application/json'
					},
					body: JSON.stringify({ image_data: img.data })
				});
			}

			const tote = await response.json();
			window.location.href = `/tote/${tote.id}`;
		} catch (error) {
			console.error('Error saving tote:', error);
			alert('Error saving tote');
		}
	} else {
		// For create mode, include all uploaded images as base64
		const toteData = {
			name,
			description,
			items,
			location,
			image_paths: uploadedImages.map(img => img.data),
			image_types: uploadedImages.map(img => img.type)
		};

		try {
			const response = await fetch('/api/tote', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify(toteData)
			});

			if (!response.ok) {
				throw new Error('Failed to create tote');
			}

			const tote = await response.json();
			window.location.href = `/tote/${tote.id}`;
		} catch (error) {
			console.error('Error creating tote:', error);
			alert('Error creating tote');
		}
	}
}

function deleteImage(imageId) {
	if (!confirm('Are you sure you want to delete this image?')) {
		return;
	}

	fetch(`/api/tote-image/${imageId}`, {
		method: 'DELETE'
	})
	.then(response => {
		if (!response.ok) {
			throw new Error('Failed to delete image');
		}
		// Reload tote data to refresh the image list
		loadToteData();
	})
	.catch(error => {
		console.error('Error deleting image:', error);
		alert('Error deleting image');
	});
}
