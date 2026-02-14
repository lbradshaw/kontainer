// form.js - Add/Edit tote form functionality

const isEditMode = window.location.pathname.includes('/edit');
let currentImagePath = null;

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
	imageInput.addEventListener('change', function(e) {
		const file = e.target.files[0];
		if (file) {
			const reader = new FileReader();
			reader.onload = function(e) {
				document.getElementById('preview-img').src = e.target.result;
				document.getElementById('image-preview').style.display = 'block';
			};
			reader.readAsDataURL(file);
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
			currentImagePath = tote.image_path;

			if (tote.image_path) {
				document.getElementById('current-image').innerHTML = 
					`<p>Current image:</p><img src="${tote.image_path}" style="max-width: 300px; max-height: 300px; border: 1px solid #ddd; border-radius: 4px;">`;
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
	const imageFile = document.getElementById('image').files[0];

	let imagePath = currentImagePath;

	// Upload image if new file selected
	if (imageFile) {
		const formData = new FormData();
		formData.append('image', imageFile);

		try {
			const response = await fetch('/api/upload-image', {
				method: 'POST',
				body: formData
			});

			if (!response.ok) {
				throw new Error('Image upload failed');
			}

			const data = await response.json();
			imagePath = data.path;
		} catch (error) {
			console.error('Error uploading image:', error);
			alert('Error uploading image');
			return;
		}
	}

	// Create or update tote
	const toteData = {
		name,
		description,
		items,
		image_path: imagePath || ''
	};

	const url = isEditMode 
		? `/api/tote/${document.getElementById('tote-id').value}`
		: '/api/tote';
	
	const method = isEditMode ? 'PUT' : 'POST';

	fetch(url, {
		method: method,
		headers: {
			'Content-Type': 'application/json'
		},
		body: JSON.stringify(toteData)
	})
	.then(response => {
		if (!response.ok) {
			throw new Error('Failed to save tote');
		}
		return response.json();
	})
	.then(tote => {
		window.location.href = `/tote/${tote.id}`;
	})
	.catch(error => {
		console.error('Error saving tote:', error);
		alert('Error saving tote');
	});
}
