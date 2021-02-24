// tailwind.config.js
module.exports = {
  purge: [
    '../templates/**/*.html',
    '../templates/**/*.js',
  ],
  darkMode: false, // or 'media' or 'class'
  theme: {
    extend: {},
  },
  variants: {
    extend: {
      opacity: ['disabled'],
    }
  },
  plugins: [],
}