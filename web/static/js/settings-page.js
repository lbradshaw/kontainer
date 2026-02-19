// settings-page.js - Settings page functionality

document.addEventListener('DOMContentLoaded', function() {
	loadSettings();
});

async function loadSettings() {
	try {
		const response = await fetch('/api/settings');
		const settings = await response.json();
		
		// Populate form fields
		document.getElementById('port').value = settings.port || 3818;
		document.getElementById('database_path').value = settings.database_path || 'kontainer.db';
		document.getElementById('theme').value = settings.theme || 'dark';
		
		// Display current config
		displayCurrentConfig(settings);
		
		// Apply theme from settings (not just localStorage)
		applyTheme(settings.theme || 'dark');
	} catch (error) {
		console.error('Error loading settings:', error);
		alert('Error loading settings');
	}
}

function displayCurrentConfig(settings) {
	const configDiv = document.getElementById('current-config');
	configDiv.innerHTML = `
		<div><strong>Port:</strong> ${settings.port || 3818}</div>
		<div><strong>Database Path:</strong> ${settings.database_path || 'kontainer.db'}</div>
		<div><strong>Theme:</strong> ${settings.theme || 'dark'}</div>
	`;
}

async function saveSettings() {
	const settings = {
		port: parseInt(document.getElementById('port').value),
		database_path: document.getElementById('database_path').value.trim(),
		theme: document.getElementById('theme').value
	};
	
	// Validate port
	if (settings.port < 1024 || settings.port > 65535) {
		alert('Port must be between 1024 and 65535');
		return;
	}
	
	// Validate database path
	if (!settings.database_path) {
		alert('Database path cannot be empty');
		return;
	}
	
	// Get current settings to check if database path changed
	let dbPathChanged = false;
	try {
		const currentResponse = await fetch('/api/settings');
		const currentSettings = await currentResponse.json();
		dbPathChanged = (currentSettings.database_path || 'kontainer.db') !== settings.database_path;
	} catch (error) {
		console.error('Error loading current settings:', error);
	}
	
	// Warn user if database path is changing
	if (dbPathChanged) {
		const confirm = window.confirm(
			'Database path is changing.\n\n' +
			'The database file will be automatically moved to the new location.\n\n' +
			'⚠️ Please ensure:\n' +
			'• The destination directory exists or can be created\n' +
			'• You have write permissions to the new location\n' +
			'• The new path is accessible\n\n' +
			'Continue with database migration?'
		);
		
		if (!confirm) {
			return;
		}
	}
	
	try {
		const response = await fetch('/api/settings', {
			method: 'PUT',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify(settings)
		});
		
		if (!response.ok) {
			const errorText = await response.text();
			throw new Error(errorText || 'Failed to save settings');
		}
		
		const savedSettings = await response.json();
		
		// Save theme to localStorage for immediate UI update
		const localSettings = JSON.parse(localStorage.getItem('kontainer_settings') || '{}');
		localSettings.theme = settings.theme;
		localStorage.setItem('kontainer_settings', JSON.stringify(localSettings));
		
		// Apply theme immediately
		applyTheme(settings.theme);
		
		// Update display
		displayCurrentConfig(savedSettings);
		
		// Show success message
		let message = 'Settings saved successfully!';
		if (dbPathChanged) {
			message += '\n\n✅ Database has been migrated to the new location.';
		}
		message += '\n\n⚠️ Please restart the application for port and database path changes to take effect.';
		
		alert(message);
	} catch (error) {
		console.error('Error saving settings:', error);
		alert('Error saving settings:\n\n' + error.message);
	}
}

function applyTheme(theme) {
	if (theme === 'dark') {
		document.documentElement.classList.add('dark-mode');
	} else {
		document.documentElement.classList.remove('dark-mode');
	}
}

async function resetSettings() {
	if (!confirm('Are you sure you want to reset all settings to defaults?')) {
		return;
	}
	
	const defaultSettings = {
		port: 3818,
		database_path: 'kontainer.db',
		theme: 'dark'
	};
	
	try {
		const response = await fetch('/api/settings', {
			method: 'PUT',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify(defaultSettings)
		});
		
		if (!response.ok) {
			throw new Error('Failed to reset settings');
		}
		
		// Clear localStorage
		localStorage.removeItem('kontainer_settings');
		
		// Reload page to show defaults
		window.location.reload();
	} catch (error) {
		console.error('Error resetting settings:', error);
		alert('Error resetting settings: ' + error.message);
	}
}

// Import data from JSON file (settings page)
async function importDataFromSettings(event) {
	const file = event.target.files[0];
	if (!file) return;

	if (!confirm(`Import data from ${file.name}?\n\nThis will ADD the totes from the file to your existing inventory.`)) {
		event.target.value = '';
		return;
	}

	try {
		const text = await file.text();
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
		window.location.href = '/';
	} catch (error) {
		console.error('Error importing data:', error);
		alert('Error importing data. Please check the file format.');
	} finally {
		event.target.value = '';
	}
}

// Delete all totes
async function deleteAllData() {
	if (!confirm('⚠️ WARNING: Delete ALL totes?\n\nThis will permanently delete all totes and their images from the database.\n\nThis action CANNOT be undone!')) {
		return;
	}

	if (!confirm('Are you ABSOLUTELY sure?\n\nType YES in the next prompt to confirm.')) {
		return;
	}

	const confirmation = prompt('Type YES to delete all data:');
	if (confirmation !== 'YES') {
		alert('Deletion cancelled.');
		return;
	}

	try {
		const response = await fetch('/api/totes/delete-all', {
			method: 'DELETE'
		});

		if (!response.ok) {
			throw new Error('Failed to delete data');
		}

		const result = await response.json();
		alert(`Successfully deleted ${result.deleted} tote(s).`);
		window.location.href = '/';
	} catch (error) {
		console.error('Error deleting data:', error);
		alert('Error deleting data.');
	}
}

