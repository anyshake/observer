# Styles

`src/index.css` is the application stylesheet entry point. It imports shared styles
before page styles, with visibility states last so hidden components stay hidden.
`loading.css` is linked directly from `index.html` to style the startup screen before
JavaScript runs.

- `shared.css`: typography, spacing, page containers, headings, and shared layouts.
- `controls.css`: buttons, form fields, and their interactive states.
- `components.css`, `layout.css`: reusable widgets and the application shell.
- `charts.css`, `markdown.css`, `syntax.css`, `notifications.css`: feature styles
  shared by multiple pages.
- `pages/`: styles specific to login, home, history, realtime, and settings.
- `tokens.css`: white and neutral gray surfaces for light mode, purple surfaces
  for dark mode, and purple accents. Includes DaisyUI theme values and colors for
  charts, maps, and MUI.
- `vendor.css`: scoped, unlayered overrides for MUI, map controls, toast surfaces,
  and DaisyUI menus.

Use existing `app-*` classes before adding new styles. Tailwind utilities are
composed with `@apply` in CSS; DaisyUI component and state classes remain in JSX.
Application styles use the utilities layer to retain their precedence over
DaisyUI defaults. Keep `vendor.css` unlayered to override MUI's unlayered styles.

Runtime dimensions remain component props/inline styles. Canvas/WebGL drawing
does not inherit CSS colors, so `useThemeColors.ts` resolves palette values for
`chartTheme.ts`, `mapTheme.ts`, and `spectrogramTheme.ts`. The MUI adapter in
`timePickerTheme.ts` reads the same tokens. Add palette values to `tokens.css`
instead of hardcoding colors in components or renderer options.

The Zustand store in `stores/theme.ts` persists the preference (`light`, `dark`,
or `system`) through the existing storage helpers. It defaults to `system`.
`selectThemeMode` resolves that preference and the current system color scheme
to a light/dark mode for renderers and MUI. `theme.ts` applies `data-theme` before
React renders and listens for system color scheme changes. Explicit light/dark
selections take precedence until the user chooses `system` again. `ThemeToggle`
is shared by the login page and application header.
Before JavaScript initializes the store, `loading.css` follows the system color
scheme for the startup screen. A saved preference takes over after initialization.

Theme changes repaint existing charts and map layers without clearing buffered
samples or recreating FFT processors. Spectrogram axes use themed colors while
retaining the original `11px AxisFont` font. Spectrogram data color maps remain
user-selectable and independent of the UI theme.
