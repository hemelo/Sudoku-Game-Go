import colors from "tailwindcss/colors";

/** @type {import('tailwindcss').Config} */
export default {
  content: ["./clients/web/**/*.templ"],
  theme: {
    extend: {
      colors: {
        primary: colors.indigo
      }
    },
  },
  plugins: [],
  safelist: [
    'space-y-1',
    'list-disc',
    'list-inside',
    'mt-2',
  ],
  darkMode: 'class'
}
