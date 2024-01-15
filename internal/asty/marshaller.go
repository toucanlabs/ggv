package asty

import (
	"go/ast"
	"go/token"
)

type AstParser struct {
	Options
	fset          *token.FileSet
	references    map[any]any
	refcount      int
	internalFuncs map[any]any
	refFuncs      map[any][]any
	externalFuncs map[any]any
}

func NewAstParser(options Options) *AstParser {
	return &AstParser{
		Options:       options,
		fset:          token.NewFileSet(),
		references:    make(map[any]any),
		internalFuncs: make(map[any]any),
		refFuncs:      make(map[any][]any),
		externalFuncs: make(map[any]any),
		refcount:      0,
	}
}

func (m *AstParser) FileSet() *token.FileSet {
	return m.fset
}

func (m *AstParser) ParseComment(comment *ast.Comment) *CommentNode {
	return &CommentNode{
		Text: comment.Text,
	}
}

func (m *AstParser) ParseComments(comments []*ast.Comment) []*CommentNode {
	if comments == nil {
		return nil
	}
	nodes := make([]*CommentNode, len(comments))
	for index, comment := range comments {
		nodes[index] = m.ParseComment(comment)
	}
	return nodes
}

func (m *AstParser) ParseCommentGroup(group *ast.CommentGroup) *CommentGroupNode {
	if !m.WithComments {
		return nil
	}
	return &CommentGroupNode{
		List: m.ParseComments(group.List),
	}
}

func (m *AstParser) ParseCommentGroups(groups []*ast.CommentGroup) []*CommentGroupNode {
	if groups == nil {
		return nil
	}
	nodes := make([]*CommentGroupNode, len(groups))
	for index, comment := range groups {
		nodes[index] = m.ParseCommentGroup(comment)
	}
	return nodes
}

// ---------------------------------------------------------------------------

func (m *AstParser) ParseField(root string, node *ast.Field) *FieldNode {
	return &FieldNode{
		Doc:     m.ParseCommentGroup(node.Doc),
		Names:   m.ParseIdents(node.Names),
		Type:    m.ParseExpr(root, node.Type),
		Tag:     m.ParseBasicLit(node.Tag),
		Comment: m.ParseCommentGroup(node.Comment),
	}
}

func (m *AstParser) ParseFields(root string, fields []*ast.Field) []*FieldNode {
	if fields == nil {
		return nil
	}
	nodes := make([]*FieldNode, len(fields))
	for index, field := range fields {
		nodes[index] = m.ParseField(root, field)
	}
	return nodes
}

func (m *AstParser) ParseFieldList(root string, node *ast.FieldList) *FieldListNode {
	return &FieldListNode{
		List: m.ParseFields(root, node.List),
	}
}

func (m *AstParser) ParseIdent(node *ast.Ident) *IdentNode {
	return &IdentNode{
		Name: node.Name,
	}
}

func (m *AstParser) ParseIdents(idents []*ast.Ident) []*IdentNode {
	if idents == nil {
		return nil
	}
	nodes := make([]*IdentNode, len(idents))
	for index, ident := range idents {
		nodes[index] = m.ParseIdent(ident)
	}
	return nodes
}

func (m *AstParser) ParseEllipsis(root string, node *ast.Ellipsis) *EllipsisNode {
	return &EllipsisNode{
		Elt: m.ParseExpr(root, node.Elt),
	}
}

func (m *AstParser) ParseBasicLit(node *ast.BasicLit) *BasicLitNode {
	return &BasicLitNode{
		Kind:  node.Kind.String(),
		Value: node.Value,
	}
}

func (m *AstParser) ParseFuncLit(root string, node *ast.FuncLit) *FuncLitNode {
	return &FuncLitNode{
		Body: m.ParseBlockStmt(root, node.Body),
	}
}

func (m *AstParser) ParseCompositeLit(root string, node *ast.CompositeLit) *CompositeLitNode {
	return &CompositeLitNode{
		Type:       m.ParseExpr(root, node.Type),
		Elts:       m.ParseExprs(root, node.Elts),
		Incomplete: node.Incomplete,
	}
}

func (m *AstParser) ParseParenExpr(root string, node *ast.ParenExpr) *ParenExprNode {
	return &ParenExprNode{
		X: m.ParseExpr(root, node.X),
	}
}

func (m *AstParser) ParseSelectorExpr(root string, node *ast.SelectorExpr) *SelectorExprNode {
	return &SelectorExprNode{
		X:   m.ParseExpr(root, node.X),
		Sel: m.ParseIdent(node.Sel),
	}
}

func (m *AstParser) ParseIndexExpr(root string, node *ast.IndexExpr) *IndexExprNode {
	return &IndexExprNode{
		X:     m.ParseExpr(root, node.X),
		Index: m.ParseExpr(root, node.Index),
	}
}

func (m *AstParser) ParseIndexListExpr(root string, node *ast.IndexListExpr) *IndexListExprNode {
	return &IndexListExprNode{
		X:       m.ParseExpr(root, node.X),
		Indices: m.ParseExprs(root, node.Indices),
	}
}

func (m *AstParser) ParseSliceExpr(root string, node *ast.SliceExpr) *SliceExprNode {
	return &SliceExprNode{
		X:      m.ParseExpr(root, node.X),
		Low:    m.ParseExpr(root, node.Low),
		High:   m.ParseExpr(root, node.High),
		Max:    m.ParseExpr(root, node.Max),
		Slice3: node.Slice3,
	}
}

func (m *AstParser) ParseTypeAssertExpr(root string, node *ast.TypeAssertExpr) *TypeAssertExprNode {
	return &TypeAssertExprNode{
		X:    m.ParseExpr(root, node.X),
		Type: m.ParseExpr(root, node.Type),
	}
}

func (m *AstParser) ParseCallExpr(root string, node *ast.CallExpr) *CallExprNode {
	return &CallExprNode{
		Fun:  m.ParseExpr(root, node.Fun),
		Args: m.ParseExprs(root, node.Args),
	}
}

func (m *AstParser) ParseStarExpr(root string, node *ast.StarExpr) *StarExprNode {
	return &StarExprNode{
		X: m.ParseExpr(root, node.X),
	}
}

func (m *AstParser) ParseUnaryExpr(root string, node *ast.UnaryExpr) *UnaryExprNode {
	return &UnaryExprNode{
		Op: node.Op.String(),
		X:  m.ParseExpr(root, node.X),
	}
}

func (m *AstParser) ParseBinaryExpr(root string, node *ast.BinaryExpr) *BinaryExprNode {
	return &BinaryExprNode{
		X:  m.ParseExpr(root, node.X),
		Op: node.Op.String(),
		Y:  m.ParseExpr(root, node.Y),
	}
}

func (m *AstParser) ParseKeyValueExpr(root string, node *ast.KeyValueExpr) *KeyValueExprNode {
	return &KeyValueExprNode{
		Key:   m.ParseExpr(root, node.Key),
		Value: m.ParseExpr(root, node.Value),
	}
}

func (m *AstParser) ParseExpr(root string, node ast.Expr) IExprNode {
	if node == nil {
		return nil
	}
	switch expr := node.(type) {
	case *ast.CallExpr:
		result := m.ParseCallExpr(root, expr)
		if _, ok := m.refFuncs[root]; ok {
			m.refFuncs[root] = append(m.refFuncs[root], result.GetFunName())
		} else {
			m.refFuncs[root] = []any{}
		}
		m.externalFuncs[result.RefId] = result
		return result
	case *ast.FuncLit:
		return m.ParseFuncLit(root, expr)
	case *ast.Ident:
		return m.ParseIdent(expr)
	case *ast.Ellipsis:
		return m.ParseEllipsis(root, expr)
	case *ast.BasicLit:
		return m.ParseBasicLit(expr)
	case *ast.CompositeLit:
		return m.ParseCompositeLit(root, expr)
	case *ast.ParenExpr:
		return m.ParseParenExpr(root, expr)
	case *ast.SelectorExpr:
		return m.ParseSelectorExpr(root, expr)
	case *ast.IndexExpr:
		return m.ParseIndexExpr(root, expr)
	case *ast.IndexListExpr:
		return m.ParseIndexListExpr(root, expr)
	case *ast.SliceExpr:
		return m.ParseSliceExpr(root, expr)
	case *ast.TypeAssertExpr:
		return m.ParseTypeAssertExpr(root, expr)
	case *ast.StarExpr:
		return m.ParseStarExpr(root, expr)
	case *ast.UnaryExpr:
		return m.ParseUnaryExpr(root, expr)
	case *ast.BinaryExpr:
		return m.ParseBinaryExpr(root, expr)
	case *ast.KeyValueExpr:
		return m.ParseKeyValueExpr(root, expr)
	case *ast.ArrayType:
		return m.ParseArrayType(root, expr)
	case *ast.StructType:
		return m.ParseStructType(root, expr)
	case *ast.InterfaceType:
		return m.ParseInterfaceType(root, expr)
	case *ast.MapType:
		return m.ParseMapType(root, expr)
	case *ast.ChanType:
		return m.ParseChanType(root, expr)
	default:
		// log.Println("implement me")
		return nil
	}
}

func (m *AstParser) ParseExprs(root string, exprs []ast.Expr) []IExprNode {
	if exprs == nil {
		return nil
	}
	nodes := []IExprNode{}
	for _, expr := range exprs {
		nodes = append(nodes, m.ParseExpr(root, expr))
	}
	return nodes
}

// ---------------------------------------------------------------------------

func (m *AstParser) ParseArrayType(root string, node *ast.ArrayType) *ArrayTypeNode {
	return &ArrayTypeNode{
		Len: m.ParseExpr(root, node.Len),
		Elt: m.ParseExpr(root, node.Elt),
	}
}

func (m *AstParser) ParseStructType(root string, node *ast.StructType) *StructTypeNode {
	return &StructTypeNode{
		Fields:     m.ParseFieldList(root, node.Fields),
		Incomplete: node.Incomplete,
	}
}

func (m *AstParser) ParseInterfaceType(root string, node *ast.InterfaceType) *InterfaceTypeNode {
	return &InterfaceTypeNode{
		Methods:    m.ParseFieldList(root, node.Methods),
		Incomplete: node.Incomplete,
	}
}

func (m *AstParser) ParseMapType(root string, node *ast.MapType) *MapTypeNode {
	return &MapTypeNode{
		Key:   m.ParseExpr(root, node.Key),
		Value: m.ParseExpr(root, node.Value),
	}
}

var ChanDirToString = map[ast.ChanDir]string{
	ast.SEND:            "SEND",
	ast.RECV:            "RECV",
	ast.SEND | ast.RECV: "BOTH",
}

func (m *AstParser) ParseChanType(root string, node *ast.ChanType) *ChanTypeNode {
	return &ChanTypeNode{
		Dir:   ChanDirToString[node.Dir],
		Value: m.ParseExpr(root, node.Value),
	}
}

// ---------------------------------------------------------------------------

func (m *AstParser) ParseBadStmt(stmt *ast.BadStmt) *BadStmtNode {
	return &BadStmtNode{}
}

func (m *AstParser) ParseDeclStmt(stmt *ast.DeclStmt) *DeclStmtNode {
	return &DeclStmtNode{
		Decl: m.ParseDecl(stmt.Decl),
	}
}

func (m *AstParser) ParseEmptyStmt(stmt *ast.EmptyStmt) *EmptyStmtNode {
	return &EmptyStmtNode{
		Implicit: stmt.Implicit,
	}
}

func (m *AstParser) ParseLabeledStmt(root string, stmt *ast.LabeledStmt) *LabeledStmtNode {
	return &LabeledStmtNode{
		Label: m.ParseIdent(stmt.Label),
		Stmt:  m.ParseStmt(root, stmt.Stmt),
	}
}

func (m *AstParser) ParseExprStmt(root string, stmt *ast.ExprStmt) *ExprStmtNode {
	return &ExprStmtNode{
		X: m.ParseExpr(root, stmt.X),
	}
}

func (m *AstParser) ParseSendStmt(root string, stmt *ast.SendStmt) *SendStmtNode {
	return &SendStmtNode{
		Chan:  m.ParseExpr(root, stmt.Chan),
		Value: m.ParseExpr(root, stmt.Value),
	}
}

func (m *AstParser) ParseIncDecStmt(root string, stmt *ast.IncDecStmt) *IncDecStmtNode {
	return &IncDecStmtNode{
		X:   m.ParseExpr(root, stmt.X),
		Tok: stmt.Tok.String(),
	}
}

func (m *AstParser) ParseAssignStmt(root string, stmt *ast.AssignStmt) *AssignStmtNode {
	return &AssignStmtNode{
		Lhs: m.ParseExprs(root, stmt.Lhs),
		Tok: stmt.Tok.String(),
		Rhs: m.ParseExprs(root, stmt.Rhs),
	}
}

func (m *AstParser) ParseGoStmt(root string, stmt *ast.GoStmt) *GoStmtNode {
	return &GoStmtNode{
		Call: m.ParseCallExpr(root, stmt.Call),
	}
}

func (m *AstParser) ParseDeferStmt(root string, stmt *ast.DeferStmt) *DeferStmtNode {
	return &DeferStmtNode{
		Call: m.ParseCallExpr(root, stmt.Call),
	}
}

func (m *AstParser) ParseReturnStmt(root string, stmt *ast.ReturnStmt) *ReturnStmtNode {
	return &ReturnStmtNode{
		Results: m.ParseExprs(root, stmt.Results),
	}
}

func (m *AstParser) ParseBranchStmt(root string, stmt *ast.BranchStmt) *BranchStmtNode {
	return &BranchStmtNode{
		Tok:   stmt.Tok.String(),
		Label: m.ParseIdent(stmt.Label),
	}
}

func (m *AstParser) ParseBlockStmt(root string, stmt *ast.BlockStmt) *BlockStmtNode {
	return &BlockStmtNode{
		List: m.ParseStmts(root, stmt.List),
	}
}

func (m *AstParser) ParseIfStmt(root string, stmt *ast.IfStmt) *IfStmtNode {
	return &IfStmtNode{
		Init: m.ParseStmt(root, stmt.Init),
		Cond: m.ParseExpr(root, stmt.Cond),
		Body: m.ParseBlockStmt(root, stmt.Body),
		Else: m.ParseStmt(root, stmt.Else),
	}
}

func (m *AstParser) ParseCaseClause(root string, stmt *ast.CaseClause) *CaseClauseNode {
	return &CaseClauseNode{
		List: m.ParseExprs(root, stmt.List),
		Body: m.ParseStmts(root, stmt.Body),
	}
}

func (m *AstParser) ParseSwitchStmt(root string, stmt *ast.SwitchStmt) *SwitchStmtNode {
	return &SwitchStmtNode{
		Init: m.ParseStmt(root, stmt.Init),
		Tag:  m.ParseExpr(root, stmt.Tag),
		Body: m.ParseBlockStmt(root, stmt.Body),
	}
}

func (m *AstParser) ParseTypeSwitchStmt(root string, stmt *ast.TypeSwitchStmt) *TypeSwitchStmtNode {
	return &TypeSwitchStmtNode{
		Init:   m.ParseStmt(root, stmt.Init),
		Assign: m.ParseStmt(root, stmt.Assign),
		Body:   m.ParseBlockStmt(root, stmt.Body),
	}
}

func (m *AstParser) ParseCommClause(root string, stmt *ast.CommClause) *CommClauseNode {
	return &CommClauseNode{
		Comm: m.ParseStmt(root, stmt.Comm),
		Body: m.ParseStmts(root, stmt.Body),
	}
}

func (m *AstParser) ParseSelectStmt(root string, stmt *ast.SelectStmt) *SelectStmtNode {
	return &SelectStmtNode{
		Body: m.ParseBlockStmt(root, stmt.Body),
	}
}

func (m *AstParser) ParseForStmt(root string, stmt *ast.ForStmt) *ForStmtNode {
	return &ForStmtNode{
		Init: m.ParseStmt(root, stmt.Init),
		Cond: m.ParseExpr(root, stmt.Cond),
		Post: m.ParseStmt(root, stmt.Post),
		Body: m.ParseBlockStmt(root, stmt.Body),
	}
}

func (m *AstParser) ParseRangeStmt(root string, stmt *ast.RangeStmt) *RangeStmtNode {
	return &RangeStmtNode{
		Key:   m.ParseExpr(root, stmt.Key),
		Value: m.ParseExpr(root, stmt.Value),
		Tok:   stmt.Tok.String(),
		X:     m.ParseExpr(root, stmt.X),
		Body:  m.ParseBlockStmt(root, stmt.Body),
	}
}

func (m *AstParser) ParseStmt(root string, node ast.Stmt) IStmtNode {
	if node == nil {
		return nil
	}
	switch stmt := node.(type) {
	case *ast.ExprStmt:
		return m.ParseExprStmt(root, stmt)
	default:
		// log.Println("implement me " + reflect.TypeOf(stmt).String())
		return nil
	}
}

func (m *AstParser) ParseStmts(root string, stmts []ast.Stmt) []IStmtNode {
	if stmts == nil {
		return nil
	}
	nodes := []IStmtNode{}
	for _, stmt := range stmts {
		if n := m.ParseStmt(root, stmt); n != nil {
			nodes = append(nodes, n)
		}
	}
	return nodes
}

func (m *AstParser) ParseFuncDecl(decl *ast.FuncDecl) *FuncDeclNode {
	node := m.ParseIdent(decl.Name)
	// internal from here
	m.internalFuncs[node.Name] = node

	return &FuncDeclNode{
		Name: node,
		Body: m.ParseBlockStmt(node.Name, decl.Body),
	}

}

func (m *AstParser) ParseDecl(node ast.Decl) IDeclNode {
	if node == nil {
		return nil
	}
	switch decl := node.(type) {
	case *ast.FuncDecl:
		return m.ParseFuncDecl(decl)
	default:
		// log.Println("implement me:", decl)
		return nil
	}
}

func (m *AstParser) ParseDecls(decls []ast.Decl) []IDeclNode {
	nodes := []IDeclNode{}
	for _, decl := range decls {
		if n := m.ParseDecl(decl); n != nil {
			nodes = append(nodes, n)
		}
	}
	return nodes
}

// ---------------------------------------------------------------------------

func (m *AstParser) ParseFile(node *ast.File) *FileNode {
	return &FileNode{
		Decls: m.ParseDecls(node.Decls),
	}
}
