import 'package:freezed_annotation/freezed_annotation.dart';
import 'package:uuid/uuid.dart';

import '../../models/label_item.dart';
import 'priority.dart';

part 'card_item.freezed.dart';

@freezed
class CardItem with _$CardItem {
  const factory CardItem({
    required String id,
    required String userId,
    required String title,
    required String description,
    required DateTime createdAt,
    required Priority? priority,
    required LabelItem? category,
    required String columnId,
  }) = _CardItem;

  factory CardItem.create({
    required String userId,
    String? title,
    String? description,
    Priority? priority,
    LabelItem? category,
    required String columnId,
  }) => CardItem(
    id: const Uuid().v4(),
    userId: userId,
    title: title ?? 'Новая задача',
    description: description ?? 'Пустое описание',
    createdAt: DateTime.now(),
    priority: priority,
    category: category,
    columnId: columnId,
  );
}

extension ShortDateFormat on DateTime {
  String toShortDate() {
    const monthNames = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];
    return '$day ${monthNames[month - 1]}';
  }
}