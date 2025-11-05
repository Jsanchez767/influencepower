/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
    "./**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        'chicago-red': '#C8102E',
        'chicago-blue': '#B3DDF2',
      },
    },
  },
  plugins: [],
}
