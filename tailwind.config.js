/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    // Note: This should be correct. and now it fully rendered.
    'frontend/htmx/error_page_handler/*.{templ,js}',
    'frontend/htmx/site/*.{templ,js}',
    'frontend/public/assets/js/*.{templ,js}',
  ],
  theme: {
    extend: {},
  },
  plugins: [],
  darkMode: 'selector',
}
