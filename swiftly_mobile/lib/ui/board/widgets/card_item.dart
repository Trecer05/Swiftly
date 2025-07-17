import 'package:freezed_annotation/freezed_annotation.dart';
import 'package:uuid/uuid.dart';

part 'card_item.freezed.dart';

enum Priority { low, medium, high }

@freezed
class CardItem with _$CardItem {
  const factory CardItem({
    required String id,
    required String title,
    required String description,
    required DateTime createdAt,
    required Priority priority,
    required String category,
    required String columnId,
  }) = _CardItem;

  factory CardItem.create({
    String? title,
    String? description,
    Priority? priority,
    String? category,
    required String columnId,
  }) => CardItem(
    id: const Uuid().v4(),
    title: title ?? 'Новая задача',
    description: description ?? 'Пустое описание',
    createdAt: DateTime.now(),
    priority: priority ?? Priority.low,
    category: category ?? 'Нет категории',
    columnId: columnId,
  );
}

extension ShortDateFormat on DateTime {
  String toShortDate() {
    const monthNames = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];
    return '$day ${monthNames[month - 1]}';
  }
}