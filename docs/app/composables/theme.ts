export enum Themes {
  DARK = 'dark',
  LIGHT = 'light',
}

export function useTheme() {
  const themeCookie = useCookie('theme', {
    maxAge: 60 * 60 * 24 * 365,
  });
  const theme = ref(themeCookie.value);

  const change = (newTheme: Themes) => {
    theme.value = newTheme;
    themeCookie.value = newTheme;
    document.documentElement.setAttribute('data-theme', newTheme);
  };

  onBeforeMount(() => {
    const storedTheme = themeCookie.value || 'dark';
    theme.value = storedTheme || 'dark';
    document.documentElement.setAttribute('data-theme', theme.value);
  });

  return {
    current: theme,
    set: change,
  }
}