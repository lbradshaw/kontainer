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
		const imageHtml = tote.image_path 
			? `<img src="${tote.image_path}" class="tote-image" alt="${tote.name}">`
			: '';
		
		const descriptionHtml = tote.description 
			? `<div class="tote-description">${tote.description}</div>`
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
				${itemsPreview}
			</div>
		`;
	}).join('');
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
