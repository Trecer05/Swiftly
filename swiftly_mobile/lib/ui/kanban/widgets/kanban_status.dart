enum KanbanStatus {
  todo('todo', 'To do'),
  progress('progress', 'In progress'),
  done('completed', 'Completed'),
  review('review', 'In review');

  final String id;
  final String title;

  const KanbanStatus(this.id, this.title);
}