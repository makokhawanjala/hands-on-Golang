// ============================================
// Windows 11 RSVP Application - JavaScript (SPA Enhanced)
// ============================================

(function() {
  'use strict';

  // === SPA Router Module ===
  const SPARouter = {
    currentRoute: null,
    contentCache: new Map(),
    
    init: function() {
      this.setupNavigationHandlers();
      this.setupMobileMenu();
      this.updateGuestCount();
      this.handleInitialRoute();
      
      // Update guest count every 30 seconds
      setInterval(() => this.updateGuestCount(), 30000);
    },
    
    // Handle initial route on page load
    handleInitialRoute: function() {
      const path = window.location.pathname;
      this.currentRoute = path;
      this.updateActiveNavItem(path);
    },
    
    // Setup navigation click handlers
    setupNavigationHandlers: function() {
      const navLinks = document.querySelectorAll('.sidebar-item');
      
      navLinks.forEach(link => {
        link.addEventListener('click', (e) => {
          const route = link.getAttribute('data-route') || link.getAttribute('href');
          
          // If it's the current route, don't reload
          if (route === this.currentRoute) {
            e.preventDefault();
            return;
          }
          
          // Allow first page load to be normal, then intercept subsequent clicks
          if (this.contentCache.size > 0) {
            e.preventDefault();
            this.navigate(route);
          }
          // On first click, let it load normally but cache the result
        });
      });
      
      // Handle browser back/forward
      window.addEventListener('popstate', () => {
        const path = window.location.pathname;
        if (this.contentCache.has(path)) {
          this.loadCachedRoute(path);
        }
      });
    },
    
    // Navigate to a route
    navigate: function(route) {
      if (route === this.currentRoute) return;
      
      this.showLoader();
      
      // If we have cached content, use it
      if (this.contentCache.has(route)) {
        this.loadCachedRoute(route);
        return;
      }
      
      // Otherwise fetch it
      this.fetchRoute(route);
    },
    
    // Fetch route content
    fetchRoute: function(route) {
      fetch(route)
        .then(response => response.text())
        .then(html => {
          // Parse the HTML and extract the content
          const parser = new DOMParser();
          const doc = parser.parseFromString(html, 'text/html');
          const content = doc.querySelector('.content-wrapper');
          
          if (content) {
            // Cache the content
            this.contentCache.set(route, content.innerHTML);
            
            // Load it
            this.loadRoute(route, content.innerHTML);
          }
        })
        .catch(error => {
          console.error('Error fetching route:', error);
          this.hideLoader();
        });
    },
    
    // Load cached route
    loadCachedRoute: function(route) {
      const content = this.contentCache.get(route);
      this.loadRoute(route, content);
    },
    
    // Load route content
    loadRoute: function(route, content) {
      const mainContent = document.getElementById('mainContent');
      const contentWrapper = mainContent.querySelector('.content-wrapper');
      
      // Fade out
      contentWrapper.classList.add('view-transition-out');
      
      setTimeout(() => {
        // Update content
        contentWrapper.innerHTML = content;
        
        // Update route
        this.currentRoute = route;
        history.pushState({}, '', route);
        
        // Update active nav item
        this.updateActiveNavItem(route);
        
        // Fade in
        contentWrapper.classList.remove('view-transition-out');
        contentWrapper.classList.add('view-transition-in');
        
        // Re-initialize modules for new content
        setTimeout(() => {
          contentWrapper.classList.remove('view-transition-in');
          this.reinitializeModules();
          this.hideLoader();
        }, 300);
        
        // Close mobile menu if open
        this.closeMobileMenu();
      }, 200);
    },
    
    // Update active navigation item
    updateActiveNavItem: function(route) {
      const navItems = document.querySelectorAll('.sidebar-item');
      navItems.forEach(item => {
        const itemRoute = item.getAttribute('data-route') || item.getAttribute('href');
        if (itemRoute === route) {
          item.classList.add('active');
        } else {
          item.classList.remove('active');
        }
      });
    },
    
    // Reinitialize modules after content change
    reinitializeModules: function() {
      FormEnhancer.init();
      TableSearch.init();
      AnimationModule.init();
      MobileEnhancements.init();
    },
    
    // Show loading indicator
    showLoader: function() {
      const loader = document.getElementById('pageLoader');
      if (loader) {
        loader.classList.add('loading');
      }
    },
    
    // Hide loading indicator
    hideLoader: function() {
      const loader = document.getElementById('pageLoader');
      if (loader) {
        setTimeout(() => {
          loader.classList.remove('loading');
        }, 300);
      }
    },
    
    // Setup mobile menu
    setupMobileMenu: function() {
      const menuToggle = document.getElementById('mobileMenuToggle');
      const sidebar = document.getElementById('sidebar');
      const overlay = document.getElementById('sidebarOverlay');
      
      if (menuToggle) {
        menuToggle.addEventListener('click', () => {
          sidebar.classList.toggle('active');
          overlay.classList.toggle('active');
          menuToggle.classList.toggle('active');
        });
      }
      
      if (overlay) {
        overlay.addEventListener('click', () => {
          this.closeMobileMenu();
        });
      }
    },
    
    // Close mobile menu
    closeMobileMenu: function() {
      const sidebar = document.getElementById('sidebar');
      const overlay = document.getElementById('sidebarOverlay');
      const menuToggle = document.getElementById('mobileMenuToggle');
      
      if (sidebar) sidebar.classList.remove('active');
      if (overlay) overlay.classList.remove('active');
      if (menuToggle) menuToggle.classList.remove('active');
    },
    
    // Update guest count in sidebar
    updateGuestCount: function() {
      fetch('/list')
        .then(response => response.text())
        .then(html => {
          const parser = new DOMParser();
          const doc = parser.parseFromString(html, 'text/html');
          const tableRows = doc.querySelectorAll('.table tbody tr');
          
          // Count actual guest rows (not empty state or no-results)
          let count = 0;
          tableRows.forEach(row => {
            if (!row.querySelector('td[colspan]')) {
              count++;
            }
          });
          
          const badge = document.getElementById('sidebarGuestCount');
          if (badge) {
            const span = badge.querySelector('span');
            if (span) {
              span.textContent = `${count} ${count === 1 ? 'Guest' : 'Guests'}`;
            }
          }
        })
        .catch(error => {
          console.error('Error updating guest count:', error);
        });
    }
  };

  // === Form Validation Module ===
  const FormValidator = {
    emailRegex: /^[^\s@]+@[^\s@]+\.[^\s@]+$/,
    phoneRegex: /^[\d\s\-\+\(\)]+$/,
    
    validateEmail: function(email) {
      if (!email || email.trim() === '') {
        return { valid: false, message: 'Email is required' };
      }
      if (!this.emailRegex.test(email)) {
        return { valid: false, message: 'Please enter a valid email address' };
      }
      return { valid: true, message: '' };
    },
    
    validatePhone: function(phone) {
      if (!phone || phone.trim() === '') {
        return { valid: false, message: 'Phone number is required' };
      }
      if (!this.phoneRegex.test(phone) || phone.replace(/\D/g, '').length < 10) {
        return { valid: false, message: 'Please enter a valid phone number (at least 10 digits)' };
      }
      return { valid: true, message: '' };
    },
    
    validateName: function(name) {
      if (!name || name.trim() === '') {
        return { valid: false, message: 'Name is required' };
      }
      if (name.trim().length < 2) {
        return { valid: false, message: 'Name must be at least 2 characters' };
      }
      return { valid: true, message: '' };
    }
  };

  // === Form Enhancement Module ===
  const FormEnhancer = {
    init: function() {
      const form = document.querySelector('form[method="POST"]');
      if (!form) return;
      
      this.form = form;
      this.setupRealTimeValidation();
      this.setupFormSubmit();
      this.enhanceInputs();
    },
    
    setupRealTimeValidation: function() {
      const nameInput = document.querySelector('input[name="name"]');
      const emailInput = document.querySelector('input[name="email"]');
      const phoneInput = document.querySelector('input[name="phone"]');
      
      if (nameInput) {
        nameInput.addEventListener('blur', () => this.validateField(nameInput, 'name'));
        nameInput.addEventListener('input', () => this.clearFieldError(nameInput));
      }
      
      if (emailInput) {
        emailInput.addEventListener('blur', () => this.validateField(emailInput, 'email'));
        emailInput.addEventListener('input', () => this.clearFieldError(emailInput));
      }
      
      if (phoneInput) {
        phoneInput.addEventListener('blur', () => this.validateField(phoneInput, 'phone'));
        phoneInput.addEventListener('input', () => this.clearFieldError(phoneInput));
      }
    },
    
    validateField: function(input, type) {
      const value = input.value;
      let result;
      
      switch(type) {
        case 'name':
          result = FormValidator.validateName(value);
          break;
        case 'email':
          result = FormValidator.validateEmail(value);
          break;
        case 'phone':
          result = FormValidator.validatePhone(value);
          break;
      }
      
      if (!result.valid) {
        this.showFieldError(input, result.message);
        return false;
      } else {
        this.showFieldSuccess(input);
        return true;
      }
    },
    
    showFieldError: function(input, message) {
      input.classList.remove('success');
      input.classList.add('error');
      
      let errorEl = input.parentElement.querySelector('.field-error');
      if (!errorEl) {
        errorEl = document.createElement('div');
        errorEl.className = 'field-error';
        errorEl.style.cssText = 'color: var(--win11-error); font-size: var(--font-size-sm); margin-top: var(--space-2);';
        input.parentElement.appendChild(errorEl);
      }
      errorEl.textContent = message;
    },
    
    showFieldSuccess: function(input) {
      input.classList.remove('error');
      input.classList.add('success');
      
      const errorEl = input.parentElement.querySelector('.field-error');
      if (errorEl) {
        errorEl.remove();
      }
    },
    
    clearFieldError: function(input) {
      input.classList.remove('error');
      const errorEl = input.parentElement.querySelector('.field-error');
      if (errorEl) {
        errorEl.remove();
      }
    },
    
    setupFormSubmit: function() {
      this.form.addEventListener('submit', (e) => {
        const nameInput = document.querySelector('input[name="name"]');
        const emailInput = document.querySelector('input[name="email"]');
        const phoneInput = document.querySelector('input[name="phone"]');
        
        let isValid = true;
        
        if (nameInput && !this.validateField(nameInput, 'name')) {
          isValid = false;
        }
        if (emailInput && !this.validateField(emailInput, 'email')) {
          isValid = false;
        }
        if (phoneInput && !this.validateField(phoneInput, 'phone')) {
          isValid = false;
        }
        
        if (!isValid) {
          e.preventDefault();
          this.scrollToFirstError();
          return false;
        }
        
        this.showLoadingState();
        
        // After form submission, update guest count
        setTimeout(() => {
          SPARouter.updateGuestCount();
        }, 1000);
      });
    },
    
    scrollToFirstError: function() {
      const firstError = document.querySelector('.form-control.error');
      if (firstError) {
        firstError.scrollIntoView({ behavior: 'smooth', block: 'center' });
        firstError.focus();
      }
    },
    
    showLoadingState: function() {
      const submitBtn = this.form.querySelector('button[type="submit"]');
      if (submitBtn) {
        submitBtn.disabled = true;
        const originalText = submitBtn.textContent;
        submitBtn.innerHTML = originalText + ' <span class="spinner"></span>';
      }
    },
    
    enhanceInputs: function() {
      const inputs = this.form.querySelectorAll('input[type="text"], input.form-control');
      inputs.forEach(input => {
        if (!input.placeholder) {
          const label = input.parentElement.querySelector('label');
          if (label) {
            input.placeholder = ' ';
          }
        }
      });
    }
  };

  // === Table Search Module ===
  const TableSearch = {
    init: function() {
      const table = document.querySelector('.table');
      if (!table) return;
      
      this.table = table;
      this.tbody = table.querySelector('tbody');
      this.rows = Array.from(this.tbody.querySelectorAll('tr'));
      
      // Don't initialize search if it's an empty state
      if (this.rows.length === 0 || (this.rows.length === 1 && this.rows[0].querySelector('td[colspan]'))) {
        return;
      }
      
      this.createSearchBar();
      this.createStats();
    },
    
    createSearchBar: function() {
      const container = this.table.parentElement;
      
      // Remove existing search if present
      const existingSearch = container.querySelector('.search-container');
      if (existingSearch) {
        existingSearch.remove();
      }
      
      const searchDiv = document.createElement('div');
      searchDiv.className = 'search-container';
      searchDiv.innerHTML = `
        <input 
          type="text" 
          class="search-input" 
          placeholder="Search guests by name, email, or phone..."
          id="guestSearch"
        />
      `;
      
      container.insertBefore(searchDiv, this.table);
      
      const searchInput = document.getElementById('guestSearch');
      searchInput.addEventListener('input', (e) => this.filterTable(e.target.value));
    },
    
    createStats: function() {
      const container = this.table.parentElement;
      
      // Remove existing stats if present
      const existingStats = container.querySelector('.stats-badge');
      if (existingStats) {
        existingStats.remove();
      }
      
      const statsDiv = document.createElement('div');
      statsDiv.className = 'stats-badge';
      statsDiv.id = 'guestStats';
      
      this.updateStats();
      
      container.insertBefore(statsDiv, container.firstChild);
    },
    
    updateStats: function() {
      const statsEl = document.getElementById('guestStats');
      if (!statsEl) return;
      
      const visibleRows = this.rows.filter(row => row.style.display !== 'none');
      const total = this.rows.length;
      const showing = visibleRows.length;
      
      if (showing === total) {
        statsEl.textContent = `${total} ${total === 1 ? 'Guest' : 'Guests'} Attending`;
      } else {
        statsEl.textContent = `Showing ${showing} of ${total} ${total === 1 ? 'Guest' : 'Guests'}`;
      }
    },
    
    filterTable: function(searchTerm) {
      const term = searchTerm.toLowerCase().trim();
      
      this.rows.forEach(row => {
        const text = row.textContent.toLowerCase();
        if (text.includes(term)) {
          row.style.display = '';
          row.style.animation = 'fadeIn 0.3s ease-in';
        } else {
          row.style.display = 'none';
        }
      });
      
      this.updateStats();
      this.showNoResults(term);
    },
    
    showNoResults: function(term) {
      let noResultsRow = this.tbody.querySelector('.no-results-row');
      
      const visibleRows = this.rows.filter(row => row.style.display !== 'none');
      
      if (visibleRows.length === 0 && term !== '') {
        if (!noResultsRow) {
          noResultsRow = document.createElement('tr');
          noResultsRow.className = 'no-results-row';
          noResultsRow.innerHTML = `
            <td colspan="3" style="text-align: center; padding: var(--space-8); color: var(--win11-text-secondary);">
              No guests found matching "${term}"
            </td>
          `;
          this.tbody.appendChild(noResultsRow);
        }
      } else {
        if (noResultsRow) {
          noResultsRow.remove();
        }
      }
    }
  };

  // === Animation Module ===
  const AnimationModule = {
    init: function() {
      this.animateOnLoad();
      this.setupHoverEffects();
    },
    
    animateOnLoad: function() {
      const cards = document.querySelectorAll('.win11-card');
      cards.forEach((card, index) => {
        card.style.opacity = '0';
        card.style.transform = 'translateY(20px)';
        
        setTimeout(() => {
          card.style.transition = 'all 0.4s ease-out';
          card.style.opacity = '1';
          card.style.transform = 'translateY(0)';
        }, index * 100);
      });
    },
    
    setupHoverEffects: function() {
      const buttons = document.querySelectorAll('.btn');
      buttons.forEach(button => {
        button.addEventListener('mouseenter', function() {
          this.style.transform = 'translateY(-2px)';
        });
        
        button.addEventListener('mouseleave', function() {
          this.style.transform = 'translateY(0)';
        });
      });
    }
  };

  // === Mobile Enhancements ===
  const MobileEnhancements = {
    init: function() {
      this.enhanceTableForMobile();
    },
    
    enhanceTableForMobile: function() {
      const table = document.querySelector('.table');
      if (!table) return;
      
      const headers = Array.from(table.querySelectorAll('thead th')).map(th => th.textContent);
      const rows = table.querySelectorAll('tbody tr');
      
      rows.forEach(row => {
        const cells = row.querySelectorAll('td');
        cells.forEach((cell, index) => {
          if (headers[index]) {
            cell.setAttribute('data-label', headers[index]);
          }
        });
      });
    }
  };

  // === Initialize Everything ===
  function init() {
    if (document.readyState === 'loading') {
      document.addEventListener('DOMContentLoaded', init);
      return;
    }
    
    // Initialize SPA Router first
    SPARouter.init();
    
    // Initialize other modules
    FormEnhancer.init();
    TableSearch.init();
    AnimationModule.init();
    MobileEnhancements.init();
    
    // Keyboard shortcuts
    document.addEventListener('keydown', function(e) {
      if ((e.ctrlKey || e.metaKey) && e.key === 'k') {
        e.preventDefault();
        const searchInput = document.getElementById('guestSearch');
        if (searchInput) {
          searchInput.focus();
        }
      }
    });
    
    console.log('Windows 11 RSVP SPA initialized successfully');
  }

  init();

  // Add CSS animations dynamically
  const style = document.createElement('style');
  style.textContent = `
    @keyframes slideInRight {
      from {
        transform: translateX(100%);
        opacity: 0;
      }
      to {
        transform: translateX(0);
        opacity: 1;
      }
    }
    
    @keyframes slideOutRight {
      from {
        transform: translateX(0);
        opacity: 1;
      }
      to {
        transform: translateX(100%);
        opacity: 0;
      }
    }
  `;
  document.head.appendChild(style);

})();