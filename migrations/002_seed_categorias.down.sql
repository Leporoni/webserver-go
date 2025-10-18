-- Remove seeded default categories
DELETE FROM categorias 
WHERE nome IN (
    'Eletrônicos',
    'Acessórios',
    'Vestuário',
    'Alimentos',
    'Casa e Jardim'
);
