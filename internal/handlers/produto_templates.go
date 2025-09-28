package handlers

// getProdutoListTemplate retorna o template da listagem de produtos
func (h *ProdutoHandler) getProdutoListTemplate() string {
	return `<!DOCTYPE html>
<html lang="pt-BR">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Produtos - Sistema de Gestão</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; background-color: #f5f5f5; }
        .container { max-width: 1400px; margin: 0 auto; background: white; padding: 20px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        .header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px; }
        .btn { padding: 10px 20px; background: #007bff; color: white; text-decoration: none; border-radius: 4px; border: none; cursor: pointer; }
        .btn:hover { background: #0056b3; }
        .btn-danger { background: #dc3545; }
        .btn-danger:hover { background: #c82333; }
        .btn-sm { padding: 5px 10px; font-size: 12px; }
        .btn-secondary { background: #6c757d; }
        .btn-secondary:hover { background: #545b62; }
        .search-box { margin-bottom: 20px; display: flex; gap: 10px; align-items: center; flex-wrap: wrap; }
        .search-box input, .search-box select { padding: 8px; border: 1px solid #ddd; border-radius: 4px; }
        .search-box input { width: 300px; }
        .search-box select { width: 200px; }
        table { width: 100%; border-collapse: collapse; margin-bottom: 20px; }
        th, td { padding: 12px; text-align: left; border-bottom: 1px solid #ddd; }
        th { background-color: #f8f9fa; font-weight: bold; }
        tr:hover { background-color: #f5f5f5; }
        .status { padding: 4px 8px; border-radius: 4px; font-size: 12px; }
        .status.ativo { background: #d4edda; color: #155724; }
        .status.inativo { background: #f8d7da; color: #721c24; }
        .pagination { display: flex; justify-content: center; gap: 10px; }
        .pagination a { padding: 8px 12px; text-decoration: none; border: 1px solid #ddd; border-radius: 4px; }
        .pagination a.current { background: #007bff; color: white; }
        .actions { display: flex; gap: 5px; }
        .price { font-weight: bold; color: #28a745; }
        .category { background: #e9ecef; padding: 2px 6px; border-radius: 3px; font-size: 11px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <div style="display: flex; align-items: center; gap: 15px;">
                <a href="/" class="btn btn-secondary">🏠 Home</a>
                <h1 style="margin: 0;">Gestão de Produtos</h1>
            </div>
            <a href="/produtos/novo" class="btn">+ Novo Produto</a>
        </div>

        <div class="search-box">
            <form method="GET" style="display: flex; gap: 10px; align-items: center; flex-wrap: wrap;">
                <input type="text" name="search" placeholder="Buscar por nome ou código..." value="{{.Search}}">
                <select name="categoria_id">
                    <option value="">Todas as categorias</option>
                    {{range .Categorias}}
                    <option value="{{.ID}}" {{if eq .ID $.CategoriaID}}selected{{end}}>{{.Nome}}</option>
                    {{end}}
                </select>
                <button type="submit" class="btn">Buscar</button>
                {{if or .Search .CategoriaID}}<a href="/produtos" class="btn btn-secondary">Limpar</a>{{end}}
            </form>
        </div>

        <p>Total: {{.Total}} produto(s)</p>

        <table>
            <thead>
                <tr>
                    <th>Nome</th>
                    <th>Código</th>
                    <th>Categoria</th>
                    <th>Preço Venda</th>
                    <th>Unidade</th>
                    <th>Status</th>
                    <th>Ações</th>
                </tr>
            </thead>
            <tbody>
                {{range .Produtos}}
                <tr>
                    <td>{{.Nome}}</td>
                    <td>{{if .CodigoBarras}}{{.CodigoBarras}}{{else}}-{{end}}</td>
                    <td>
                        {{if .Categoria}}
                            <span class="category">{{.Categoria.Nome}}</span>
                        {{else}}-{{end}}
                    </td>
                    <td class="price">{{formatPrice .PrecoVenda}}</td>
                    <td>{{.UnidadeMedida}}</td>
                    <td>
                        <span class="status {{if .Ativo}}ativo{{else}}inativo{{end}}">
                            {{if .Ativo}}Ativo{{else}}Inativo{{end}}
                        </span>
                    </td>
                    <td class="actions">
                        <a href="/produtos/{{.ID}}" class="btn btn-sm">Ver</a>
                        <a href="/produtos/{{.ID}}/editar" class="btn btn-sm">Editar</a>
                        <button onclick="deleteProduto('{{.ID}}', '{{.Nome}}')" class="btn btn-sm btn-danger">Excluir</button>
                    </td>
                </tr>
                {{else}}
                <tr>
                    <td colspan="7" style="text-align: center; padding: 40px;">
                        {{if or .Search .CategoriaID}}
                            Nenhum produto encontrado para os filtros aplicados
                        {{else}}
                            Nenhum produto cadastrado
                        {{end}}
                    </td>
                </tr>
                {{end}}
            </tbody>
        </table>

        {{if gt .TotalPages 1}}
        <div class="pagination">
            {{if .HasPrev}}
                <a href="?page={{sub .Page 1}}{{if .Search}}&search={{.Search}}{{end}}{{if .CategoriaID}}&categoria_id={{.CategoriaID}}{{end}}">&laquo; Anterior</a>
            {{end}}
            
            {{range $i := .TotalPages}}
                {{$page := add $i 1}}
                <a href="?page={{$page}}{{if $.Search}}&search={{$.Search}}{{end}}{{if $.CategoriaID}}&categoria_id={{$.CategoriaID}}{{end}}" 
                   {{if eq $page $.Page}}class="current"{{end}}>{{$page}}</a>
            {{end}}
            
            {{if .HasNext}}
                <a href="?page={{add .Page 1}}{{if .Search}}&search={{.Search}}{{end}}{{if .CategoriaID}}&categoria_id={{.CategoriaID}}{{end}}">Próximo &raquo;</a>
            {{end}}
        </div>
        {{end}}
    </div>

    <script>
        function deleteProduto(id, nome) {
            if (confirm('Tem certeza que deseja excluir o produto "' + nome + '"?')) {
                fetch('/api/produtos/' + id, {
                    method: 'DELETE'
                })
                .then(response => response.json())
                .then(data => {
                    if (data.success) {
                        alert('Produto excluído com sucesso!');
                        location.reload();
                    } else {
                        alert('Erro ao excluir produto: ' + data.error);
                    }
                })
                .catch(error => {
                    alert('Erro ao excluir produto: ' + error);
                });
            }
        }
    </script>
</body>
</html>`
}

// getProdutoDetailTemplate retorna o template de detalhes do produto
func (h *ProdutoHandler) getProdutoDetailTemplate() string {
	return `<!DOCTYPE html>
<html lang="pt-BR">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Nome}} - Produto</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; background-color: #f5f5f5; }
        .container { max-width: 1000px; margin: 0 auto; background: white; padding: 20px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        .header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px; }
        .btn { padding: 10px 20px; background: #007bff; color: white; text-decoration: none; border-radius: 4px; }
        .btn:hover { background: #0056b3; }
        .btn-secondary { background: #6c757d; }
        .btn-secondary:hover { background: #545b62; }
        .info-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 30px; margin-bottom: 20px; }
        .info-section { background: #f8f9fa; padding: 20px; border-radius: 8px; }
        .info-item { margin-bottom: 15px; }
        .info-label { font-weight: bold; color: #666; margin-bottom: 5px; }
        .info-value { font-size: 16px; }
        .status { padding: 4px 8px; border-radius: 4px; font-size: 14px; }
        .status.ativo { background: #d4edda; color: #155724; }
        .status.inativo { background: #f8d7da; color: #721c24; }
        .actions { display: flex; gap: 10px; margin-top: 20px; }
        .price { font-size: 18px; font-weight: bold; color: #28a745; }
        .category { background: #007bff; color: white; padding: 4px 8px; border-radius: 4px; font-size: 12px; }
        .image-section { text-align: center; }
        .product-image { max-width: 200px; max-height: 200px; border-radius: 8px; border: 1px solid #ddd; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <div style="display: flex; align-items: center; gap: 15px;">
                <a href="/" class="btn btn-secondary">🏠 Home</a>
                <h1 style="margin: 0;">{{.Nome}}</h1>
            </div>
            <span class="status {{if .Ativo}}ativo{{else}}inativo{{end}}">
                {{if .Ativo}}Ativo{{else}}Inativo{{end}}
            </span>
        </div>

        <div class="info-grid">
            <div class="info-section">
                <h3>Informações Básicas</h3>
                <div class="info-item">
                    <div class="info-label">Código de Barras:</div>
                    <div class="info-value">{{if .CodigoBarras}}{{.CodigoBarras}}{{else}}-{{end}}</div>
                </div>
                <div class="info-item">
                    <div class="info-label">Categoria:</div>
                    <div class="info-value">
                        {{if .Categoria}}
                            <span class="category">{{.Categoria.Nome}}</span>
                        {{else}}-{{end}}
                    </div>
                </div>
                <div class="info-item">
                    <div class="info-label">Unidade de Medida:</div>
                    <div class="info-value">{{.UnidadeMedida}}</div>
                </div>
                <div class="info-item">
                    <div class="info-label">Peso:</div>
                    <div class="info-value">{{if .Peso}}{{.Peso}} kg{{else}}-{{end}}</div>
                </div>
                <div class="info-item">
                    <div class="info-label">Dimensões:</div>
                    <div class="info-value">{{if .Dimensoes}}{{.Dimensoes}}{{else}}-{{end}}</div>
                </div>
            </div>

            <div class="info-section">
                <h3>Preços e Margem</h3>
                <div class="info-item">
                    <div class="info-label">Preço de Custo:</div>
                    <div class="info-value">{{formatPricePtr .PrecoCusto}}</div>
                </div>
                <div class="info-item">
                    <div class="info-label">Preço de Venda:</div>
                    <div class="info-value price">{{formatPrice .PrecoVenda}}</div>
                </div>
                <div class="info-item">
                    <div class="info-label">Margem de Lucro:</div>
                    <div class="info-value">{{if .MargemLucro}}{{.MargemLucro}}%{{else}}-{{end}}</div>
                </div>
            </div>
        </div>

        {{if .Descricao}}
        <div class="info-section">
            <h3>Descrição</h3>
            <p>{{.Descricao}}</p>
        </div>
        {{end}}

        {{if .ImagemURL}}
        <div class="info-section image-section">
            <h3>Imagem</h3>
            <img src="{{.ImagemURL}}" alt="{{.Nome}}" class="product-image">
        </div>
        {{end}}

        {{if .Observacoes}}
        <div class="info-section">
            <h3>Observações</h3>
            <p>{{.Observacoes}}</p>
        </div>
        {{end}}

        <div class="info-section">
            <div class="info-item">
                <div class="info-label">Cadastrado em:</div>
                <div class="info-value">{{.CreatedAt.Format "02/01/2006 15:04"}}</div>
            </div>
            <div class="info-item">
                <div class="info-label">Última atualização:</div>
                <div class="info-value">{{.UpdatedAt.Format "02/01/2006 15:04"}}</div>
            </div>
        </div>

        <div class="actions">
            <a href="/produtos" class="btn btn-secondary">← Voltar</a>
            <a href="/produtos/{{.ID}}/editar" class="btn">Editar</a>
        </div>
    </div>
</body>
</html>`
}

// getProdutoFormTemplate retorna o template do formulário de produto
func (h *ProdutoHandler) getProdutoFormTemplate() string {
	return `<!DOCTYPE html>
<html lang="pt-BR">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}} - Sistema de Gestão</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; background-color: #f5f5f5; }
        .container { max-width: 1000px; margin: 0 auto; background: white; padding: 20px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        .header { margin-bottom: 20px; }
        .form-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 20px; }
        .form-group { margin-bottom: 15px; }
        .form-group.full-width { grid-column: 1 / -1; }
        label { display: block; margin-bottom: 5px; font-weight: bold; color: #333; }
        input, select, textarea { width: 100%; padding: 8px; border: 1px solid #ddd; border-radius: 4px; font-size: 14px; }
        textarea { height: 80px; resize: vertical; }
        .checkbox-group { display: flex; align-items: center; gap: 10px; }
        .checkbox-group input { width: auto; }
        .btn { padding: 10px 20px; background: #007bff; color: white; text-decoration: none; border-radius: 4px; border: none; cursor: pointer; }
        .btn:hover { background: #0056b3; }
        .btn-secondary { background: #6c757d; }
        .btn-secondary:hover { background: #545b62; }
        .actions { display: flex; gap: 10px; margin-top: 20px; }
        .error { background: #f8d7da; color: #721c24; padding: 10px; border-radius: 4px; margin-bottom: 20px; }
        .required { color: red; }
        .form-section { background: #f8f9fa; padding: 15px; border-radius: 8px; margin-bottom: 20px; }
        .form-section h3 { margin-top: 0; color: #333; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <div style="display: flex; align-items: center; gap: 15px;">
                <a href="/" class="btn btn-secondary">🏠 Home</a>
                <h1 style="margin: 0;">{{.Title}}</h1>
            </div>
        </div>

        {{if .Error}}
        <div class="error">{{.Error}}</div>
        {{end}}

        <form method="POST" action="{{.Action}}">
            <div class="form-section">
                <h3>Informações Básicas</h3>
                <div class="form-grid">
                    <div class="form-group">
                        <label for="nome">Nome <span class="required">*</span></label>
                        <input type="text" id="nome" name="nome" required 
                               value="{{if .Produto}}{{.Produto.Nome}}{{end}}">
                    </div>

                    <div class="form-group">
                        <label for="codigo_barras">Código de Barras</label>
                        <input type="text" id="codigo_barras" name="codigo_barras" 
                               value="{{if and .Produto .Produto.CodigoBarras}}{{.Produto.CodigoBarras}}{{end}}">
                    </div>

                    <div class="form-group">
                        <label for="categoria_id">Categoria</label>
                        <select id="categoria_id" name="categoria_id">
                            <option value="">Selecione uma categoria</option>
                            {{range .Categorias}}
                            <option value="{{.ID}}" {{if and $.Produto $.Produto.CategoriaID (eq .ID $.Produto.CategoriaID)}}selected{{end}}>{{.Nome}}</option>
                            {{end}}
                        </select>
                    </div>

                    <div class="form-group">
                        <label for="unidade_medida">Unidade de Medida <span class="required">*</span></label>
                        <select id="unidade_medida" name="unidade_medida" required>
                            <option value="UN" {{if and .Produto (eq .Produto.UnidadeMedida "UN")}}selected{{end}}>Unidade (UN)</option>
                            <option value="KG" {{if and .Produto (eq .Produto.UnidadeMedida "KG")}}selected{{end}}>Quilograma (KG)</option>
                            <option value="G" {{if and .Produto (eq .Produto.UnidadeMedida "G")}}selected{{end}}>Grama (G)</option>
                            <option value="L" {{if and .Produto (eq .Produto.UnidadeMedida "L")}}selected{{end}}>Litro (L)</option>
                            <option value="ML" {{if and .Produto (eq .Produto.UnidadeMedida "ML")}}selected{{end}}>Mililitro (ML)</option>
                            <option value="M" {{if and .Produto (eq .Produto.UnidadeMedida "M")}}selected{{end}}>Metro (M)</option>
                            <option value="CM" {{if and .Produto (eq .Produto.UnidadeMedida "CM")}}selected{{end}}>Centímetro (CM)</option>
                        </select>
                    </div>
                </div>

                <div class="form-group full-width">
                    <label for="descricao">Descrição</label>
                    <textarea id="descricao" name="descricao">{{if and .Produto .Produto.Descricao}}{{.Produto.Descricao}}{{end}}</textarea>
                </div>
            </div>

            <div class="form-section">
                <h3>Preços e Margem</h3>
                <div class="form-grid">
                    <div class="form-group">
                        <label for="preco_custo">Preço de Custo (R$)</label>
                        <input type="number" id="preco_custo" name="preco_custo" step="0.01" min="0"
                               value="{{if and .Produto .Produto.PrecoCusto}}{{formatFloat .Produto.PrecoCusto}}{{end}}">
                    </div>

                    <div class="form-group">
                        <label for="preco_venda">Preço de Venda (R$) <span class="required">*</span></label>
                        <input type="number" id="preco_venda" name="preco_venda" step="0.01" min="0.01" required
                               value="{{if .Produto}}{{.Produto.PrecoVenda}}{{end}}">
                    </div>

                    <div class="form-group">
                        <label for="margem_lucro">Margem de Lucro (%)</label>
                        <input type="number" id="margem_lucro" name="margem_lucro" step="0.01" min="0" max="100"
                               value="{{if and .Produto .Produto.MargemLucro}}{{formatFloat .Produto.MargemLucro}}{{end}}">
                    </div>
                </div>
            </div>

            <div class="form-section">
                <h3>Características Físicas</h3>
                <div class="form-grid">
                    <div class="form-group">
                        <label for="peso">Peso (kg)</label>
                        <input type="number" id="peso" name="peso" step="0.001" min="0"
                               value="{{if and .Produto .Produto.Peso}}{{formatFloat .Produto.Peso}}{{end}}">
                    </div>

                    <div class="form-group">
                        <label for="dimensoes">Dimensões (ex: 10x20x30 cm)</label>
                        <input type="text" id="dimensoes" name="dimensoes"
                               value="{{if and .Produto .Produto.Dimensoes}}{{.Produto.Dimensoes}}{{end}}">
                    </div>

                    <div class="form-group full-width">
                        <label for="imagem_url">URL da Imagem</label>
                        <input type="url" id="imagem_url" name="imagem_url"
                               value="{{if and .Produto .Produto.ImagemURL}}{{.Produto.ImagemURL}}{{end}}">
                    </div>
                </div>
            </div>

            <div class="form-section">
                <h3>Observações e Status</h3>
                <div class="form-group">
                    <label for="observacoes">Observações</label>
                    <textarea id="observacoes" name="observacoes">{{if and .Produto .Produto.Observacoes}}{{.Produto.Observacoes}}{{end}}</textarea>
                </div>

                <div class="form-group">
                    <div class="checkbox-group">
                        <input type="checkbox" id="ativo" name="ativo" 
                               {{if not .Produto}}checked{{else}}{{if .Produto.Ativo}}checked{{end}}{{end}}>
                        <label for="ativo">Produto Ativo</label>
                    </div>
                </div>
            </div>

            <div class="actions">
                <a href="/produtos" class="btn btn-secondary">Cancelar</a>
                <button type="submit" class="btn">{{if .IsEdit}}Atualizar{{else}}Criar{{end}} Produto</button>
            </div>
        </form>
    </div>

    <script>
        // Calcular margem de lucro automaticamente
        document.getElementById('preco_custo').addEventListener('input', calcularMargem);
        document.getElementById('preco_venda').addEventListener('input', calcularMargem);

        function calcularMargem() {
            const precoCusto = parseFloat(document.getElementById('preco_custo').value) || 0;
            const precoVenda = parseFloat(document.getElementById('preco_venda').value) || 0;
            
            if (precoCusto > 0 && precoVenda > 0) {
                const margem = ((precoVenda - precoCusto) / precoCusto) * 100;
                document.getElementById('margem_lucro').value = margem.toFixed(2);
            }
        }
    </script>
</body>
</html>`
}