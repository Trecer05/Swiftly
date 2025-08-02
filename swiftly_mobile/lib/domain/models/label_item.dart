import 'dart:ui';

import 'package:freezed_annotation/freezed_annotation.dart';
import 'package:uuid/uuid.dart';

part 'label_item.freezed.dart';

@freezed
class LabelItem with _$LabelItem{
  const factory LabelItem({
    required String id,
    required String? cardId,
    required String? userId,
    required String title,
    required Color color,
  }) = _LabelItem;

  factory LabelItem.create({
    String? cardId,
    String? userId,
    required String title,
    required Color color,
  }) {
    assert((cardId == null) != (userId == null), 'Exactly one of cardId or userId must be non-null.');
    return LabelItem(
    id: const Uuid().v4(),
    cardId: cardId,
    userId: userId,
    title: title,
    color: color,
  );
  }
}